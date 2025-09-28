namespace go cp_center

struct CPMaterial {
    1: i64 MaterialID
    2: i64 CpID
    3: string CpIcon
    4: string CpName
    5: list<string> VerificationImages
    6: string BusinessLicenses
    7: string Website
    8: MaterialStatus Status
    9: string ReviewComment
    10: i64 CreateTime
    11: i64 ModifyTime
}

enum MaterialStatus {
    Unset = 0
    Draft = 1
    Reviewing = 2
    Online = 3
    Rejected = 4
}

enum SubmitMode {
    Unset = 0
    SubmitDraft = 1
    SubmitReview = 2
}

struct CreateCPMaterialsRequest {
    1: CPMaterial CPMaterial
    2: SubmitMode SubmitMode
}

struct CreateCPMaterialResponse {
    // 255: BaseResp BaseResp
}

service CpCenterService {

}