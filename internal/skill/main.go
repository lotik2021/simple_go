package skill

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"github.com/Masterminds/semver"
)

const (
	ObjectID = "movista_skill"
)

func FindForVersion(ctx common.Context, v string) (skills []Skill, err error) {

	skills = make([]Skill, 0)

	err = ctx.DB.Model(&skills).Order("ord asc").Select()
	if err != nil {
		return
	}

	if v == "" {
		return
	}

	ver, err := semver.NewVersion(v)
	if err != nil {
		return
	}

	filteredSkills := make([]Skill, 0, len(skills))

	for _, val := range skills {
		if val.DisableInVersions == "" {
			filteredSkills = append(filteredSkills, val)
			continue
		}

		constraints, err := semver.NewConstraint(val.DisableInVersions)
		if err != nil {
			logger.Log.Errorf("cannot create version constraints for skill %+v", v)
			continue
		}

		result := constraints.Check(ver)

		if !result {
			filteredSkills = append(filteredSkills, val)
		}
	}

	return filteredSkills, err
}
