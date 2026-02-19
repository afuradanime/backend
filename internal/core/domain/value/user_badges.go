package value

type UserBadges uint8

const (
	UserBadgeContributor       UserBadges = iota + 1 // Has contributted to the project
	UserBadgeTranslator                              // Has contributed in translations to the site
	UserBadgeBrand                                   // Is a brand account
	UserBadgeBetaTester                              // Was a beta tester
	UserBadgeSuperMegaIllyaFan                       // Illyasviel von Einzbern  (イリヤスフィール・フォン・アインツベルン, Iriyasufīru fon Aintsuberun?), often referred to as Illya[Note 2] (イリヤ, Iriya?), is the Master of Berserker in the Fifth Holy Grail War of Fate/stay night.
)
