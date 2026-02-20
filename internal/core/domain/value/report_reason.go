package value

type ReportReason uint8

const (
	ReportReasonHateSpeech ReportReason = iota + 1
	ReportReasonViolenceOrHarassment
	ReportReasonPrivacy
	ReportReasonBot
	ReportReasonSpam
	ReportReasonIllegalActivities
	ReportReasonMisinformation
	ReportReasonOther
)
