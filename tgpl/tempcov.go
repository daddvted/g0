package tgpl

import "fmt"

type Celsius float64
type Fahrenheig float64

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsiu  = 100
)

func (c Celsius) String() string { return fmt.Sprintf }
