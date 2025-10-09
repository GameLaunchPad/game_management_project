namespace go game_platform_api
/*
   api服务和rpc服务差别:
   1. id类需要定义为string类型
   2. 结构体需要按照下划线命名
   3. 新接口需要带上path
   ref:https://www.cloudwego.io/zh/docs/kitex/tutorials/advanced-feature/generic-call/thrift_idl_annotation_standards/
*/

include "common.thrift"

// content provider
struct CreateCPMaterialsRequest {
    1: CPMaterial cp_material
    2: SubmitMode submit_mode
}

struct CreateCPMaterialResponse {
    1: CreateCPMaterialData data
    255: common.BaseResp base_resp
}

struct CreateCPMaterialData {
    1: string cp_id
    2: string material_id
}

struct CPMaterial {
    1: string material_id
    2: string cp_id
    3: string cp_icon
    4: string cp_name
    5: list<string> verification_images
    6: string business_license
    7: string website
    8: MaterialStatus status
    9: string review_comment
    10: i64 create_time
    11: i64 modify_time
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

struct UpdateCPMaterialsRequest {
    1: i64 material_id (api.path = 'id')
    2: CPMaterial cp_material
    3: SubmitMode submit_mode
}

struct UpdateCPMaterialData {

}

struct UpdateCPMaterialResponse {
    255: common.BaseResp base_resp
}

struct ReviewCPMaterialRequest {
    1: i64 material_id
    2: i64 cp_id
    3: ReviewResult review_result
    4: ReviewRemark review_remark
}

struct ReviewRemark {
    1: string remark
    2: string operator
    3: i64 review_time
    4: string meta
}

struct ReviewCPMaterialData {

}

struct ReviewCPMaterialResponse {
    1: ReviewCPMaterialData data
    255: common.BaseResp base_resp
}

enum ReviewResult {
    Unset = 0
    Pass = 1
    Reject = 2
}

struct GetCPMaterialRequest {
    1: string material_id (api.path = 'id')
    2: string cp_id
}

struct GetCPMaterialResponse {
    1: GetCPMaterialData data
    255: common.BaseResp base_resp
}

struct GetCPMaterialData {
    1: CPMaterial cp_material
}

// Game
struct GetGameListRequest {
   1: optional GameListFilter filter
   2: optional GameListSorter sorter
   3: i32 page_num
   4: i32 page_size
}

struct GameListFilter {
    1: optional string filter_text
}

struct GameListSorter {
    1: optional i64 update_time
}

struct GetGameListResponse {
    1: GetGameListData data
    255: common.BaseResp base_resp
}

struct GetGameListData {
    1: list<BriefGame> game_list
    2: i32 total_count
}

struct BriefGame {
    1: string game_id
    2: string cp_id
    3: string game_name
    4: string game_icon
    5: i64 create_time
    6: i64 update_time
    7: GameStatus game_status
}

enum GameStatus {
    Unset = 0
    Draft = 1
    Reviewing = 2
    Published = 3
    Rejected = 4
}


struct GetGameDetailRequest {
   1: i64 game_id (api.path = 'id')
}

struct GetGameDetailResponse {
    1: GetGameDetailData data
    255: common.BaseResp base_resp
}

struct GetGameDetailData {
    1: GameDetail game_detail
}

struct GameDetail {
    1: string game_id
    2: string cp_id
    3: GameVersion online_game_version
    4: GameVersion newest_game_version
    5: i64 create_time
    6: i64 modify_time
}

struct GameVersion {
    1: string game_id
    2: string game_version_id
    3: string game_name
    4: string game_icon
    5: string game_introduction
    6: list<string> game_introduction_images
    7: string header_image
    8: list<GamePlatform> game_platforms
    9: string package_name
    10: string download_url
    11: GameStatus game_status
    12: ReviewRemark review_remark
    13: i64 create_time
    14: i64 update_time
}

enum GamePlatform {
    Unset = 0
    Android = 1
    IOS = 2
    Web = 3
}

struct GameDetailWrite {
    1: string game_id
    2: i64 cp_id
    3: GameVersion game_version
}

struct CreateGameDetailRequest {
   1: GameDetailWrite game_detail
   2: SubmitMode submit_mode
}

struct CreateGameDetailResponse {
    1:  CreateGameDetailData data
    255: common.BaseResp BaseResp
}

struct CreateGameDetailData {
    1: string game_id
}

struct UpdateGameDetailRequest {
    1: string game_id(api.path = 'id')
    2: GameDetailWrite game_detail
    3: SubmitMode submit_mode
}

struct UpdateGameDetailResponse {
    1: UpdateGameDetailData data
    2: common.BaseResp base_resp
}

struct UpdateGameDetailData {

}

struct ReviewGameVersionRequest {
    1: string game_id
    2: string game_version_id
    3: ReviewResult review_result
    4: ReviewRemark review_remark
}

struct ReviewGameVersionResponse {
    1: ReviewGameVersionData data
    2: common.BaseResp base_resp
}

struct ReviewGameVersionData {
}

service GamePlatformAPIService {
     // content provider
     CreateCPMaterialResponse CreateCPMaterial(1: CreateCPMaterialsRequest req) (api.post = '/api/v1/cp/materials') // 创建厂商材料
     UpdateCPMaterialResponse UpdateCPMaterial (1: UpdateCPMaterialsRequest req) (api.put = '/api/v1/cp/materials/:id') // 更新厂商材料
     ReviewCPMaterialResponse ReviewCPMaterial (1: ReviewCPMaterialRequest req) (api.post = '/api/v1/cp/materials/review') // 审核厂商材料
     GetCPMaterialResponse GetCPMaterial(1: GetCPMaterialRequest req) (api.get = '/api/v1/cp/materials/:id') // 获取厂商材料

     // games management
     GetGameListResponse GetGameList(1: GetGameListRequest req) (api.get = '/api/v1/games') // 获取游戏列表
     GetGameDetailResponse GetGameDetail(1: GetGameDetailRequest req) (api.get = '/api/v1/games/:id') // 获取游戏详情
     CreateGameDetailResponse CreateGameDetail (1: CreateGameDetailRequest req) (api.post = '/api/v1/games') // 创建游戏信息
     UpdateGameDetailResponse UpdateGameDetail(1: UpdateGameDetailRequest req) (api.put = '/api/v1/games/:id') // 更新游戏信息
     ReviewGameVersionResponse ReviewGameVersion(1: ReviewGameVersionRequest req) (api.post = '/api/v1/games/review') // 审核游戏信息
}