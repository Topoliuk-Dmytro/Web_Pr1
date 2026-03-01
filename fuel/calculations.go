package fuel

import (
	"errors"
	"math"
)

// WorkingMassComposition - склад робочої маси палива (%)
type WorkingMassComposition struct {
	HP float64 // водень
	CP float64 // вуглець
	SP float64 // сірка
	NP float64 // азот
	OP float64 // кисень
	WP float64 // волога
	AP float64 // зола
}

// Task1Result - результат Завдання 1
type Task1Result struct {
	KRS           float64  // коефіцієнт переходу робоча -> суха
	KRG           float64  // коефіцієнт переходу робоча -> горюча
	DryMass       DryMassComposition
	Combustible   CombustibleMassComposition
	QrH           float64 // нижча теплота згоряння робочої маси, МДж/кг
	QdH           float64 // нижча теплота згоряння сухої маси, МДж/кг
	QgH           float64 // нижча теплота згоряння горючої маси, МДж/кг
}

// DryMassComposition - склад сухої маси
type DryMassComposition struct {
	Hd, Cd, Sd, Nd, Od, Ad float64
}

// CombustibleMassComposition - склад горючої маси
type CombustibleMassComposition struct {
	Hg, Cg, Sg, Ng, Og float64
}

// Task1Calculate - Завдання 1: розрахунок складу сухої та горючої маси та теплоти згоряння
func Task1Calculate(w WorkingMassComposition) (Task1Result, error) {
	if w.WP >= 100 {
		return Task1Result{}, errors.New("вологість Wₚ не може бути >= 100%")
	}
	if w.WP+w.AP >= 100 {
		return Task1Result{}, errors.New("сума вологості та золи не може бути >= 100%")
	}
	// Коефіцієнти переходу (табл. 1.1)
	krs := 100 / (100 - w.WP)
	krg := 100 / (100 - w.WP - w.AP)

	// Склад сухої маси (множник з табл. 1.1: робоча -> суха = 100/(100-Wr))
	dry := DryMassComposition{
		Hd: round(w.HP * krs),
		Cd: round(w.CP * krs),
		Sd: round(w.SP * krs),
		Nd: round(w.NP * krs),
		Od: round(w.OP * krs),
		Ad: round(w.AP * krs),
	}

	// Склад горючої маси
	comb := CombustibleMassComposition{
		Hg: round(w.HP * krg),
		Cg: round(w.CP * krg),
		Sg: round(w.SP * krg),
		Ng: round(w.NP * krg),
		Og: round(w.OP * krg),
	}

	// Нижча теплота згоряння за формулою Менделєєва (1.2)
	// QРН = 339СР + 1030НР - 108,8(ОР - SР) - 25WР, кДж/кг
	qrhKj := 339*w.CP + 1030*w.HP - 108.8*(w.OP-w.SP) - 25*w.WP
	qrh := qrhKj / 1000 // МДж/кг

	// Перерахунок теплоти згоряння (табл. 1.2)
	// Qdi = Qri * 100/(100-Wr)
	// Qdafi = Qri * 100/(100-Wr-Ar)
	qdh := qrh * 100 / (100 - w.WP)
	qgh := qrh * 100 / (100 - w.WP - w.AP)

	return Task1Result{
		KRS:         round(krs),
		KRG:         round(krg),
		DryMass:     dry,
		Combustible: comb,
		QrH:         round(qrh),
		QdH:         round(qdh),
		QgH:         round(qgh),
	}, nil
}

// CombustibleMassInput - вхідні дані для Завдання 2 (склад горючої маси мазуту)
type CombustibleMassInput struct {
	Cg    float64 // вуглець, %
	Hg    float64 // водень, %
	Og    float64 // кисень, %
	Sg    float64 // сірка, %
	QgH   float64 // нижча теплота згоряння горючої маси, МДж/кг
	Wr    float64 // вологість робочої маси, %
	Ad    float64 // зольність сухої маси, %
	Vg    float64 // вміст ванадію, мг/кг
}

// Task2Result - результат Завдання 2
type Task2Result struct {
	WorkingMass WorkingMassComposition
	Vr          float64 // ванадій на робочу масу, мг/кг
	QrH         float64 // нижча теплота згоряння робочої маси, МДж/кг
}

// Task2Calculate - Завдання 2: перерахунок з горючої маси на робочу
func Task2Calculate(in CombustibleMassInput) (Task2Result, error) {
	if in.Wr >= 100 {
		return Task2Result{}, errors.New("вологість не може бути >= 100%")
	}
	// Ar = Ad * (100 - Wr) / 100
	ar := in.Ad * (100 - in.Wr) / 100

	// Множник переходу горюча -> робоча: (100 - Wr - Ar) / 100 (табл. 1.1)
	multCombToWork := (100 - in.Wr - ar) / 100
	// Множник для золи: суха -> робоча: (100 - Wr) / 100
	multDryToWork := (100 - in.Wr) / 100

	working := WorkingMassComposition{
		CP: round(in.Cg * multCombToWork),
		HP: round(in.Hg * multCombToWork),
		OP: round(in.Og * multCombToWork),
		SP: round(in.Sg * multCombToWork),
		NP: 0, // азот не задається в завданні 2
		WP: in.Wr,
		AP: round(in.Ad * multDryToWork),
	}

	// Ванадій: Vr = Vg * (100 - Wr) / 100
	vr := in.Vg * multDryToWork

	// Теплота згоряння: Qri = Qdafi * (100 - Wr - Ar) / 100 (табл. 1.2)
	qrh := in.QgH * (100 - in.Wr - ar) / 100

	return Task2Result{
		WorkingMass: working,
		Vr:          round(vr),
		QrH:         round(qrh),
	}, nil
}

func round(x float64) float64 {
	return math.Round(x*100) / 100
}
