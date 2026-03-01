package fuel

import (
	"math"
	"testing"
)

func TestTask1_ControlExample(t *testing.T) {
	// Контрольний приклад з PDF: HP=1,9%; CP=21,1%; SP=2,60%; NP=0,20%; OP=7,10%; WP=53,0%; AP=14,1
	w := WorkingMassComposition{HP: 1.9, CP: 21.1, SP: 2.6, NP: 0.2, OP: 7.1, WP: 53, AP: 14.1}
	r, err := Task1Calculate(w)
	if err != nil {
		t.Fatal(err)
	}
	// Очікувані: KRS=2.13, KRG=3.04
	if math.Abs(r.KRS-2.13) > 0.01 {
		t.Errorf("KRS = %v, очікувалось 2.13", r.KRS)
	}
	if math.Abs(r.KRG-3.04) > 0.01 {
		t.Errorf("KRG = %v, очікувалось 3.04", r.KRG)
	}
	// QrH = 7.2953 МДж/кг
	if math.Abs(r.QrH-7.2953) > 0.001 {
		t.Errorf("QrH = %v, очікувалось 7.2953", r.QrH)
	}
}

func TestTask2_ControlExample(t *testing.T) {
	// Контрольний приклад: C=85.5, H=11.2, O=0.8, S=2.5, Q=40.4, W=2, Ad=0.15, V=333.3
	in := CombustibleMassInput{Cg: 85.5, Hg: 11.2, Og: 0.8, Sg: 2.5, QgH: 40.4, Wr: 2, Ad: 0.15, Vg: 333.3}
	r, err := Task2Calculate(in)
	if err != nil {
		t.Fatal(err)
	}
	// Очікувані: CP=83.66, HP=10.96, SP=2.45, OP=0.78, AP=0.15, Vr=326.63, QrH=39.48
	if math.Abs(r.WorkingMass.CP-83.66) > 0.1 {
		t.Errorf("CP = %v, очікувалось 83.66", r.WorkingMass.CP)
	}
	if math.Abs(r.WorkingMass.HP-10.96) > 0.1 {
		t.Errorf("HP = %v, очікувалось 10.96", r.WorkingMass.HP)
	}
	if math.Abs(r.QrH-39.48) > 0.1 {
		t.Errorf("QrH = %v, очікувалось 39.48", r.QrH)
	}
}
