namespace go game

struct GetGameListRequest {
   1: optional GameListFilter Filter
   2: optional GameListSorter Sorter
   3: i32 PageNum
   4: i32 PageSize
}

struct GameListFilter {
    1: optional string FilterText
}

struct GameListSorter {
    1: optional i64 UpdateTime
}

struct GetGameListResponse {
    1: list<BriefGame> GameList
    2: i32 TotalCount
    //255: BaseResp BaseResp
}

struct BriefGame {
    1: i64 GameID
    2: i64 CpID
    3: string GameName
    4: string GameIcon
    5: i64 CreateTime
    6: i64 UpdateTime
    7: GameStatus GameStatus
}

enum GameStatus {
    Unset = 0
    Draft = 1 // 草稿态
    Reviewing = 2 // 审核中
    Published = 3 // 已发布
    Rejected = 4 // 已拒绝
}

struct GetGameDetailRequest {
   1: i64 GameID
}

struct GetGameDetailResponse {
    1: GameDetail GameDetail
    // 255: BaseResp BaseResp
}

struct GameDetail {
    1: i64 GameID
    2: i64 CpID
    3: GameVersion OnlineGameVersion
    4: GameVersion NewestGameVersion // 创建新游戏时仅需填写NewestGameVersion字段
    5: i64 CreateTime
    6: i64 ModifyTime
}

struct GameVersion {
    1: i64 GameID
    2: i64 GamVersionID
    3: string GameName
    4: string GameIcon
    5: string GameIntroduction
    6: list<string> GameIntroductionImages
    7: list<GamePlatform> GamePlatforms
    8: string PackageName // 包名
    9: string DownloadURL // 下载链接
    10: GameStatus GameStatus // 游戏状态
    11: string ReviewComment
    12: i64 ReviewTime
    13: i64 CreateTime
    14: i64 UpdateTime
}

enum GamePlatform {
    Unset = 0
    Android = 1
    IOS = 2
    Web = 3
}

// 更新时使用，做读写区分的原因是字段较大差异
struct GameDetailWrite {
    1: i64 GameID
    2: i64 CpID
    3: GameVersion GameVersion
}

struct CreateGameDetailRequest {
   1: GameDetailWrite GameDetail
   2: SubmitMode SubmitMode
}

struct CreateGameDetailResponse {
    1: i64 GameID
    // 255: BaseResp base_resp
}


enum SubmitMode {
    Unset = 0
    SubmitDraft = 1 // 提交草稿
    SubmitReview = 2 // 提交审核
}

struct ReviewGameVersionRequest {
   1: i64 GameID
   2: i64 GameVersionID
   3: ReviewResult review_result
}

struct ReviewGameVersionResponse {
    //255: BaseResp BaseResp
}

enum ReviewResult {
    Unset = 0
    Pass = 1
    Reject = 2
}

service GameService {
    GetGameListResponse GetGameList (1: GetGameListRequest req) // 获取游戏列表
    GetGameDetailResponse GetGameDetail (1: GetGameDetailRequest req) // 获取游戏详情
    CreateGameDetailResponse CreateGameDetail (1: CreateGameDetailRequest req) // 创建游戏详情
    ReviewGameVersionResponse ReviewGameVersion (1: ReviewGameVersionRequest req) // 审核游戏信息
}

