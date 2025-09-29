namespace go cp_center

include "common.thrift"

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

struct CreateCPMaterialRequest {
    1: CPMaterial CPMaterial
    2: SubmitMode SubmitMode
}

struct CreateCPMaterialResponse {
    1: i64 CpID
    2: i64 MaterialID
    255: common.BaseResp BaseResp
}

struct UpdateCPMaterialRequest {
    1: i64 MaterialID
    2: CPMaterial CPMaterial
    3: SubmitMode SubmitMode
}

struct UpdateCPMaterialResponse {
   255: common.BaseResp BaseResp
}

struct ReviewCPMaterialRequest {
    1: i64 CpID
    2: i64 MaterialID
    3: ReviewResult review_result
}

enum ReviewResult {
    Unset = 0
    Pass = 1
    Reject = 2
}

struct ReviewCPMaterialResponse {
    255: common.BaseResp BaseResp
}

struct GetCPMaterialRequest {
    1: i64 CpID
}

struct GetCPMaterialResponse {
    1: CPMaterial CPMaterial
    255: common.BaseResp BaseResp
}

service CpCenterService {
    CreateCPMaterialResponse CreateCPMaterial (1: CreateCPMaterialRequest req) // 创建认证材料
    UpdateCPMaterialResponse UpdateCPMaterial (1: UpdateCPMaterialRequest req) // 更新认证材料
    ReviewCPMaterialResponse ReviewCPMaterial (1: ReviewCPMaterialRequest req) // 审核厂商材料
    GetCPMaterialResponse GetCPMaterial(1: GetCPMaterialRequest req) // 获取厂商认证材料
}