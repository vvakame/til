/* tslint:disable */
/* eslint-disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL subscription operation: CommentAdded
// ====================================================

export interface CommentAdded_commentAdded {
  __typename: "Comment";
  id: string;
  text: string;
}

export interface CommentAdded {
  commentAdded: CommentAdded_commentAdded | null;
}
