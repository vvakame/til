export interface AEBackupInformationDeleteReq {
    key?: string;
}
export interface ApiCheckedCircleOwn {
    cursor?: string;
    eventExhibitCourceID?: string;
    eventID?: string;
    limit?: number;
    offset?: number;
    userID?: string;
    visibility?: "site" | "staff";
}
export interface ApiCircle {
    cursor?: string;
    eventExhibitCourceID?: string;
    eventID?: string;
    limit?: number;
    offset?: number;
    onlyAdoption?: boolean;
    visibility?: "site" | "staff";
}
export interface ApiCircleOwn {
    visibility?: "site" | "staff";
}
export interface ApiCircleTicket {
    circleExhibitInfoID?: number;
    cursor?: string;
    eventID?: string;
    limit?: number;
    offset?: number;
    userEmail?: string;
    userID?: number;
    visibility?: "site" | "staff";
}
export interface ApiEvent {
    cursor?: string;
    limit?: number;
    offset?: number;
    visibility?: "site" | "staff";
}
export interface ApiMarketHandshake {
    buyerID?: number;
    cursor?: string;
    distributorID?: number;
    eventID?: string;
    limit?: number;
    offset?: number;
    status?: "concluded" | "canceled";
    visibility: "buyer" | "distributor" | "staff";
}
export interface ApiMigrateFirebaseDbMarkethandshakeUpdateDb {
    cursor?: string;
    limit?: number;
    offset?: number;
}
export interface ApiProduct {
    circleExhibitInfoID?: string;
    cursor?: string;
    eventID?: string;
    limit?: number;
    offset?: number;
    visibility?: "site" | "own" | "staff";
}
export interface ApiProductcontent {
    cursor?: string;
    limit?: number;
    offset?: number;
    productInfoID?: string;
    visibility?: "site" | "own" | "staff";
}
export interface ApiStaffCircle {
    cursor?: string;
    eventExhibitCourceID?: string;
    eventID?: string;
    limit?: number;
    offset?: number;
    onlyAdoption?: boolean;
    visibility?: "site" | "staff";
}
export interface ApiUser {
    cursor?: string;
    limit?: number;
    offset?: number;
}
export interface ApiUserChangeEmailCallback {
    token?: string;
}
export interface ApiUserSignupCallback {
    token?: string;
}
export interface CheckedCircleExhibit {
    circleExhibitInfo?: CircleExhibitInfo;
    note?: string;
    type?: "manual" | "derived";
    userID?: string; // int64
}
export interface CheckedCircleListResp {
    cursor?: string;
    list?: CheckedCircleExhibit[];
}
export interface CircleExhibitInfo {
    agreements?: boolean;
    checkedCount?: number; // int32
    circleCutImage?: Image;
    copyFromCircleExhibitInfoID?: string; // int64
    createdAt?: string; // date-time
    event?: Event;
    eventExhibitCourse?: EventExhibitCourse;
    genre?: "software" | "hardware" | "technology" | "other";
    genreFreeFormat?: string;
    grayscaleCircleCutImage?: Image;
    id?: string; // int64
    name?: string;
    nameRuby?: string;
    nextCircleExhibitInfoID?: string; // int64
    note?: string;
    numberOfPublication?: number; // int32
    penName?: string;
    prevCircleExhibitInfoID?: string; // int64
    spaces?: string[];
    staffNote?: string;
    stockOfPublication?: number; // int32
    updatedAt?: string; // date-time
    userIDs?: number /* int64 */[];
    webSiteURL?: string;
}
export interface CircleExhibitInfoSecret {
    address?: string;
    bankAccountNameKana?: string;
    bankAccountNumber?: string;
    bankAccountType?: string;
    bankBranchID?: string;
    bankBranchName?: string;
    bankID?: string;
    bankName?: string;
    email?: string;
    marketPaymentType?: string;
    paymentAgreements?: boolean;
    phoneNumber?: string;
    postalCode?: string;
    representativeName?: string;
    representativeNameRuby?: string;
    requiredTickets?: number; // int32
    status?: string;
}
export interface CircleExhibitInfoWithSecret {
    address?: string;
    agreements?: boolean;
    bankAccountNameKana?: string;
    bankAccountNumber?: string;
    bankAccountType?: string;
    bankBranchID?: string;
    bankBranchName?: string;
    bankID?: string;
    bankName?: string;
    checkedCount?: number; // int32
    circleCutImage?: Image;
    copyFromCircleExhibitInfoID?: string; // int64
    createdAt?: string; // date-time
    email?: string;
    event?: Event;
    eventExhibitCourse?: EventExhibitCourse;
    genre?: "software" | "hardware" | "technology" | "other";
    genreFreeFormat?: string;
    grayscaleCircleCutImage?: Image;
    id?: string; // int64
    marketPaymentType?: string;
    name?: string;
    nameRuby?: string;
    nextCircleExhibitInfoID?: string; // int64
    note?: string;
    numberOfPublication?: number; // int32
    paymentAgreements?: boolean;
    penName?: string;
    phoneNumber?: string;
    postalCode?: string;
    prevCircleExhibitInfoID?: string; // int64
    representativeName?: string;
    representativeNameRuby?: string;
    requiredTickets?: number; // int32
    spaces?: string[];
    staffNote?: string;
    status?: string;
    stockOfPublication?: number; // int32
    updatedAt?: string; // date-time
    userIDs?: number /* int64 */[];
    webSiteURL?: string;
}
export interface CircleListResp {
    cursor?: string;
    list?: CircleExhibitInfo[];
}
export interface CirclePaypalCreatePaymentReq {
    circleExhibitInfoID?: string; // int64
    eventID?: string;
}
export interface CirclePaypalCreatePaymentResp {
    id?: string;
    paymentTransactionID?: string; // int64
}
export interface CirclePaypalExecutePaymentReq {
    circleExhibitInfoID?: string; // int64
    eventID?: string;
    payerID?: string;
    paymentID?: string;
    paymentTransactionID?: string; // int64
}
export interface CircleSpacePatch {
    circleID?: string; // int64
    circleName?: string;
    genre?: "software" | "hardware" | "technology" | "other";
    spaces?: string[];
    staffNote?: string;
}
export interface CircleSpacesPatchReq {
    dryRun?: boolean;
    eventID?: string;
    list?: CircleSpacePatch[];
}
export interface CircleStaffListResp {
    cursor?: string;
    list?: CircleExhibitInfoWithSecret[];
}
export interface CircleStripeCreateChargeReq {
    circleExhibitInfoID?: string; // int64
    eventID?: string;
    token?: string;
}
export interface CircleTicket {
    circleExhibitInfo?: CircleExhibitInfo;
    createdAt?: string; // date-time
    id?: string; // int64
    suspended?: boolean;
    updatedAt?: string; // date-time
    usedAtList?: string /* date-time */[];
    user?: User;
    userEmail?: string;
    userID?: string; // int64
}
export interface CircleTicketAcceptReq {
    id?: string; // int64
    usedAtList?: string /* date-time */[];
}
export interface CircleTicketListResp {
    cursor?: string;
    list?: CircleTicket[];
}
export interface CircleTicketStaffDistributeReq {
    cursor?: string;
    eventID?: string;
}
export interface CircleTicketTransferReq {
    id?: string; // int64
    userEmail?: string;
    visibility?: "site" | "staff";
}
export interface Event {
    catalogImageDeadline?: string; // date-time
    createdAt?: string; // date-time
    electionClosed?: boolean;
    endAt?: string; // date-time
    eventExhibitCourses?: EventExhibitCourse[];
    id?: string;
    name?: string;
    ogpDescription?: string;
    ogpImage?: string;
    ogpTitle?: string;
    place?: string;
    recruitEndAt?: string; // date-time
    recruitStartAt?: string; // date-time
    startAt?: string; // date-time
    updatedAt?: string; // date-time
}
export interface EventExhibitCourse {
    createdAt?: string; // date-time
    endAt?: string; // date-time
    exhibitFee?: number; // int32
    id?: string; // int64
    name?: string;
    place?: string;
    startAt?: string; // date-time
    updatedAt?: string; // date-time
}
export interface EventListResp {
    cursor?: string;
    list?: Event[];
}
export interface ExportProductPaymentQRCodToStorageReq {
    bucket?: string;
    circleExhibitInfoID?: number; // int64
    size?: number; // int32
}
export interface ExportProductPaymentQRCodeReq {
    bucket?: string;
    eventID?: string;
    size?: number; // int32
}
export interface FirebaseDBMarketHandshakeCheckReq {
    operator?: string;
    userID?: string; // int64
}
export interface GCSObject {
    bucket?: string;
    contentType?: string;
    crc32c?: string;
    etag?: string;
    generation?: string;
    id?: string;
    md5Hash?: string;
    mediaLink?: string;
    metageneration?: string;
    name?: string;
    selfLink?: string;
    size?: string; // int64
    timeCreated?: string; // date-time
    timeDeleted?: string; // date-time
    updated?: string; // date-time
}
export interface GCSObjectToBQJobReq {
    TimeCreated?: string; // date-time
    bucket?: string;
    filePath?: string;
    kindName?: string;
}
export interface Image {
    createdAt?: string; // date-time
    fileSize?: number; // int32
    height?: number; // int32
    id?: string; // int64
    updatedAt?: string; // date-time
    url?: string;
    width?: number; // int32
}
export interface IntIDInPathReq {
    id?: string; // int64
}
export interface IntTokenCallbackReq {
    token?: string; // int64
}
export interface MailApplyBillingTemplateReq {
    cursor?: string;
    eventID?: string;
    limit?: number; // int32
    offset?: number; // int32
    status?: string;
    subject?: string;
    template?: string;
}
export interface MailApplyCircleTemplateReq {
    adoption?: boolean;
    cursor?: string;
    eventExhibitCourseID?: string; // int64
    eventID?: string;
    limit?: number; // int32
    offset?: number; // int32
    subject?: string;
    template?: string;
}
export interface MailApplyTemplateResp {
    cursor?: string;
    list?: MailInfo[];
}
export interface MailBatchSendReq {
    list?: MailInfo[];
}
export interface MailInfo {
    createdAt?: string; // date-time
    from?: string;
    fromName?: string;
    html?: string;
    id?: string;
    replyTo?: string;
    sGResponse?: string;
    send?: boolean;
    subject?: string;
    text?: string;
    tos?: string[];
    updatedAt?: string; // date-time
}
export interface MailSendPromotionMailReq {
    from?: string;
    fromName?: string;
    html?: string;
    replyTo?: string;
    subject?: string;
    text?: string;
}
export interface MarketBillPaypalCreatePaymentReq {
    billID?: string; // int64
}
export interface MarketBillPaypalCreatePaymentResp {
    id?: string;
    paymentTransactionID?: string; // int64
}
export interface MarketBillPaypalExecutePaymentReq {
    billID?: string; // int64
    payerID?: string;
    paymentID?: string;
    paymentTransactionID?: string; // int64
}
export interface MarketBillStripeCreateChargeReq {
    billID?: string; // int64
    token?: string;
}
export interface MarketHandshake {
    buyerAgreement?: "approve";
    buyerChargeRate?: number; // double
    buyerChargeTotal?: number; // int32
    buyerID?: string; // int64
    contentRevision?: number; // int32
    createdAt?: string; // date-time
    distributorAgreement?: "approve" | "cancel";
    distributorChargeRate?: number; // double
    distributorChargeTotal?: number; // int32
    distributorCircleExhibitInfo?: CircleExhibitInfo;
    id?: string; // int64
    items?: MarketHandshakeProductDetail[];
    lockedAt?: string; // date-time
    sharedCode?: string;
    status?: string;
    updatedAt?: string; // date-time
}
export interface MarketHandshakeBill {
    createdAt?: string; // date-time
    erasures?: boolean[];
    handshakes?: MarketHandshake[];
    id?: string; // int64
    status?: string;
    subTotals?: number /* int32 */[];
    total?: number; // int32
    updatedAt?: string; // date-time
    userID?: string; // int64
}
export interface MarketHandshakeListResp {
    cursor?: string;
    list?: MarketHandshake[];
}
export interface MarketHandshakeProductDetail {
    price?: number; // int32
    product?: ProductInfo;
    quantity?: number; // int32
}
export interface MarketHandshakeUpdateReq {
    handshake?: MarketHandshake;
    id?: string; // int64
    operator?: "buyer" | "distributor";
    visibility?: "site" | "staff";
}
export interface Noop {
}
export interface ProductContent {
    contentType?: string;
    createdAt?: string; // date-time
    fileName?: string;
    fileSize?: string; // int64
    id?: string; // int64
    productInfoID?: string; // int64
    updatedAt?: string; // date-time
}
export interface ProductContentInsertReq {
    productInfoID?: string; // int64
    tempProductContentID?: string; // int64
}
export interface ProductContentListResp {
    cursor?: string;
    list?: ProductContent[];
}
export interface ProductContentUploadContentReq {
    productInfoID?: string; // int64
}
export interface ProductContentUploadContentResp {
    uploadURL?: string;
}
export interface ProductContentUploadContentSucceedResp {
    tempProductContentID?: string; // int64
}
export interface ProductInfo {
    circleExhibitInfoID?: string; // int64
    contentRevision?: number; // int32
    createdAt?: string; // date-time
    description?: string;
    firstAppearanceEventName?: string;
    firstAtTechBookFest?: boolean;
    id?: string; // int64
    images?: Image[];
    name?: string;
    page?: number; // int32
    price?: number; // int32
    relatedURLs?: string[];
    seq?: number; // int32
    stock?: number; // int32
    training?: boolean;
    type?: "commerce" | "fanzine";
    updatedAt?: string; // date-time
}
export interface ProductInfoListResp {
    cursor?: string;
    list?: ProductInfo[];
}
export interface ProsessMailReq {
    targetIDs?: string[];
}
export interface PushInfo {
    createdAt?: string; // date-time
    id?: string;
    token?: string;
    type?: "fcm-android" | "fcm-ios";
    updatedAt?: string; // date-time
}
export interface ReqListOptions {
    cursor?: string;
    limit?: number; // int32
    offset?: number; // int32
}
export interface RespListOptions {
    cursor?: string;
}
export interface SMSSendReq {
    message?: string;
    toPhoneNumber?: string;
}
export interface Setting {
    boolValue?: boolean;
    createdAt?: string; // date-time
    floatValue?: number; // double
    id?: string;
    integerArrayValue?: number /* int64 */[];
    integerValue?: string; // int64
    stringValue?: string;
    timeValue?: string; // date-time
    updatedAt?: string; // date-time
}
export interface StaffReqOptions {
    visibility?: "site" | "staff";
}
export interface TempUser {
    agreements?: boolean;
    createdAt?: string; // date-time
    email?: string;
    id?: string;
    name?: string;
    newPassword?: string;
    updatedAt?: string; // date-time
}
export interface TqDatastoreManagementDeleteOldBackups {
    cursor?: string;
    limit?: number;
    offset?: number;
}
export interface TqMarketBillMake {
    buyerID?: number;
    cursor?: string;
    distributorID?: number;
    eventID?: string;
    limit?: number;
    offset?: number;
    status?: "concluded" | "canceled";
    visibility: "buyer" | "distributor" | "staff";
}
export interface TqMarketTransferMake {
    buyerID?: number;
    cursor?: string;
    distributorID?: number;
    eventID?: string;
    limit?: number;
    offset?: number;
    status?: "concluded" | "canceled";
    visibility: "buyer" | "distributor" | "staff";
}
export interface User {
    agreements?: boolean;
    aliases?: string[];
    createdAt?: string; // date-time
    disablePromotionMail?: boolean;
    email?: string;
    id?: string; // int64
    name?: string;
    newPassword?: string;
    phoneNumber?: string;
    staff?: boolean;
    updatedAt?: string; // date-time
}
export interface UserChangeEmailReq {
    newEmail?: string;
    password?: string;
}
export interface UserChangePhoneNumberReq {
    phoneNumber?: string;
}
export interface UserCommitChangePhoneNumberReq {
    verifyCode?: string;
}
export interface UserCommitResetPassword {
    email?: string;
    newPassword?: string;
    token?: string; // int64
}
export interface UserFirebaseCustomTokenResp {
    token?: string;
}
export interface UserListResp {
    cursor?: string;
    list?: User[];
}
export interface UserLoginRequest {
    email?: string;
    password?: string;
}
export interface UserResetPasswordReq {
    email?: string;
}
export interface UserRevokeAuthReq {
    type?: "google" | "github" | "twitter";
}
export interface UserUpdateReq {
    agreements?: boolean;
    aliases?: string[];
    createdAt?: string; // date-time
    disablePromotionMail?: boolean;
    email?: string;
    id?: string; // int64
    name?: string;
    newPassword?: string;
    oldPassword?: string;
    phoneNumber?: string;
    staff?: boolean;
    updatedAt?: string; // date-time
}
export interface UserWithdrawReq {
    password?: string;
}
