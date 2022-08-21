package consts

const (
	MsgTypeWarn = "warning"
)

const (
	Msg400BadRequest      = "不正なリクエストです"
	Msg500UnexpectedError = "予期せぬエラーが発生しました"
)

const (
	MsgNeedToLogin                 = "ログインしてください"
	MsgFailedToRegisterValidations = "カスタムバリデーションの登録に失敗しました"
	MsgFailedToLogin               = "メールアドレスまたはパスワードが違います"
	MsgDuplicatedEmailAddress      = "登録済みのメールアドレスです"
	MsgFailedToRegisterCircleIcon  = "サークルアイコンの登録に失敗しました"
)

const (
	ErrMsgFailedToLoadFile     = "ファイルの読み込みに失敗しました"
	ErrMsgFailedToRegisterToS3 = "S3への登録に失敗しました"
)
