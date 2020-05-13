package dictionary

import "bitbucket.movista.ru/maas/maasapi/internal/common"

func GetCitizenships(ctx common.Context, cultureCode string) ([]Citizenship, error) {
	citShips := []Citizenship{}

	err := ctx.DB.Model(&Citizenship{}).Where("lng = ?", cultureCode).Select(&citShips)
	if err != nil {
		return nil, err
	}

	return citShips, nil
}

func GetCitizenshipByAlpha(ctx common.Context, alphaCode string) (Citizenship, error) {
	citShips := Citizenship{}

	err := ctx.DB.Model(&Citizenship{}).Where("alpha2 = ?", alphaCode).Select(&citShips)
	if err != nil {
		return Citizenship{}, err
	}

	return citShips, nil
}
