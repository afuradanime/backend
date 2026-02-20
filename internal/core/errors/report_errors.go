package domain_errors

type CannotReportYourselfError struct {
	AnimeID string
}

func (e CannotReportYourselfError) Error() string {
	return "Cannot report yourself"
}

type AlreadyReportedError struct {
	AnimeID string
}

func (e AlreadyReportedError) Error() string {
	return "You've already reported this account"
}

type ReportNotFoundError struct {
	AnimeID string
}

func (e ReportNotFoundError) Error() string {
	return "Report not found"
}
