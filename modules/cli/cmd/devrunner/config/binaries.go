package config

// Url holds the parts of a URL that point to a release.
type Url struct {
	Prefix string
	Middle string
	Suffix string
}

// AvailableServiceReleaseUrls holds the parts of the release URLs for known services.
type AvailableServiceReleaseUrls struct {
	CometBFT        Url
	AstriaSequencer Url
	AstriaComposer  Url
	AstriaConductor Url
}

// CometBftReleaseUrl returns the release URLs for the known CometBFT service.
func (asru *AvailableServiceReleaseUrls) CometBftReleaseUrl(version string) string {
	return asru.CometBFT.Prefix + version + asru.CometBFT.Middle + version + asru.CometBFT.Suffix
}

// AstriaSequencerReleaseUrl returns the release URLs for the known Astria Sequencer service.
func (asru *AvailableServiceReleaseUrls) AstriaSequencerReleaseUrl(version string) string {
	return asru.AstriaSequencer.Prefix + version + asru.AstriaSequencer.Suffix
}

// AstriaComposerReleaseUrl returns the release URLs for the known Astria Composer service.
func (asru *AvailableServiceReleaseUrls) AstriaComposerReleaseUrl(version string) string {
	return asru.AstriaComposer.Prefix + version + asru.AstriaComposer.Suffix
}

// AstriaConductorReleaseUrl returns the release URLs for the known Astria Conductor service.
func (asru *AvailableServiceReleaseUrls) AstriaConductorReleaseUrl(version string) string {
	return asru.AstriaConductor.Prefix + version + asru.AstriaConductor.Suffix
}
