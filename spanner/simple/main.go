package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/admin/database/apiv1"
	"github.com/pborman/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type RefreshToken struct {
	RefreshTokenID string
	Token          string
	ClientId       string
	Scopes         []string
	Nonce          string
}

type AccessToken struct {
	RefreshTokenID string
	AccessTokenID  string
	Token          string
}

func main() {
	ctx := context.Background()

	db := fmt.Sprintf(
		"projects/%s/instances/%s/databases/%s",
		os.Getenv("SPANNER_PROJECT"),
		os.Getenv("SPANNER_INSTANCE"),
		os.Getenv("SPANNER_DB"),
	)
	log.Println("connect to", db)

	adminClient, err := database.NewDatabaseAdminClient(
		ctx,
		option.WithCredentialsFile("spanner-sa.json"),
	)
	if err != nil {
		log.Fatal(err)
	}

	dataClient, err := spanner.NewClient(
		ctx,
		db,
		option.WithCredentialsFile("spanner-sa.json"),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = createDatabase(ctx, adminClient, db)
	if err != nil {
		log.Println(err)
	}

	err = createSchema(ctx, adminClient, db)
	if err != nil {
		log.Println(err)
	}

	refreshToken := uuid.New()
	err = insertRefreshToken(ctx, dataClient, refreshToken)
	if err != nil {
		log.Fatalln(err)
	}

	accessToken := uuid.New()
	err = insertAccessToken(ctx, dataClient, refreshToken, accessToken)
	if err != nil {
		log.Fatalln(err)
	}

	found, err := checkAccessToken(ctx, dataClient, accessToken)
	if err != nil {
		log.Println(err)
	}
	if !found {
		log.Fatalln("access token not found")
	}

	err = revokeRefreshToken(ctx, dataClient, refreshToken)
	if err != nil {
		log.Fatalln(err)
	}
}

func createDatabase(ctx context.Context, adminClient *database.DatabaseAdminClient, db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}
	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
	})
	if err != nil {
		return err
	}
	_, err = op.Wait(ctx)
	if err != nil {
		return err
	}

	log.Printf("Created database [%s]\n", db)

	return nil
}

func createSchema(ctx context.Context, adminClient *database.DatabaseAdminClient, db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE RefreshTokens (
				RefreshTokenId	STRING(1024) NOT NULL,
				Token			STRING(1024) NOT NULL,
				ClientId		STRING(1024) NOT NULL,
				Scopes			Array<STRING(1024)>,
				Nonce			STRING(1024) NOT NULL
			) PRIMARY KEY (RefreshTokenId)`,
			`CREATE TABLE AccessTokens (
				RefreshTokenId	STRING(1024) NOT NULL,
				AccessTokenId	STRING(1024) NOT NULL,
				Token			STRING(1024) NOT NULL,
			) PRIMARY KEY (RefreshTokenId, AccessTokenId),
			INTERLEAVE IN PARENT RefreshTokens ON DELETE CASCADE`,
		},
	})
	if err != nil {
		return err
	}
	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	log.Printf("Update DDL [%s]\n", db)

	return nil
}

func insertRefreshToken(ctx context.Context, client *spanner.Client, refreshToken string) error {
	m, err := spanner.InsertStruct("RefreshTokens", &RefreshToken{
		RefreshTokenID: refreshToken,
		Token:          "Foobar",
		ClientId:       "Client1",
		Scopes: []string{
			"foobarscope",
		},
		Nonce: "nonce1",
	})
	if err != nil {
		return err
	}
	ms := []*spanner.Mutation{
		m,
	}
	_, err = client.Apply(ctx, ms)
	return err
}

func insertAccessToken(ctx context.Context, client *spanner.Client, refreshToken string, accessToken string) error {
	m, err := spanner.InsertStruct("AccessTokens", &AccessToken{
		RefreshTokenID: refreshToken,
		AccessTokenID:  accessToken,
		Token:          accessToken,
	})
	if err != nil {
		return err
	}
	ms := []*spanner.Mutation{
		m,
	}
	_, err = client.Apply(ctx, ms)
	return err
}

func revokeRefreshToken(ctx context.Context, client *spanner.Client, refreshToken string) error {
	m := []*spanner.Mutation{
		spanner.Delete("RefreshTokens", spanner.Key{refreshToken}),
	}
	_, err := client.Apply(ctx, m)
	return err
}

func checkAccessToken(ctx context.Context, client *spanner.Client, accessToken string) (bool, error) {
	stmt := spanner.Statement{
		SQL: `SELECT AccessTokenId FROM AccessTokens WHERE AccessTokenId = @accessToken`,
		Params: map[string]interface{}{
			"accessToken": accessToken,
		},
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return false, nil
		}
		if err != nil {
			return false, err
		}

		fetchedAccessToken := &AccessToken{}
		if err := row.ToStruct(fetchedAccessToken); err != nil {
			return false, err
		}
		log.Println("fetched access token", fetchedAccessToken)
		return true, nil
	}
}
