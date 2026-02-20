package domain_errors

type CannotReportYourselfError struct {
	AnimeID string
}

func (e CannotReportYourselfError) Error() string {
	return "Cannot report yourself"
}

type UserAlreadyRestrictedError struct {
	AnimeID string
}

func (e UserAlreadyRestrictedError) Error() string {
	return "This user is already restricted"
}

type AlreadyReportedError struct {
	AnimeID string
}

func (e AlreadyReportedError) Error() string {
	return "You've already reported this account"
}

type ReportNotFoundError struct{}

func (e ReportNotFoundError) Error() string {
	return "Report not found"
}
