package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	nameMaxSize       = 42
	endOfString       = 0x00
	countOfBitsInByte = 8

	// Bit masks and shifts for respectStrength
	respectMask  = 0xF0 // [1111 0000]
	respectShift = 4
	strengthMask = 0x0F // [0000 1111]

	// Bit masks and shifts for levelExp
	levelMask      = 0xF0 // [1111 0000]
	levelShift     = 4
	experienceMask = 0x0F // [0000 1111]

	// Bit masks and shifts for typeHouseGunFamily
	typeMask   = 0xC0 // [1100 0000]
	typeShift  = 6
	houseFlag  = 0x04 // [0000 0100]
	gunFlag    = 0x02 // [0000 0010]
	familyFlag = 0x01 // [0000 0001]

	// Bit masks and shifts for manaHealth
	manaFirstMask  = 0xC0 // [1100 0000] - Top 2 bits of first byte
	manaFirstShift = 6

	healthFirstMask = 0x03 // [0000 0011] - Bottom 2 bits of first byte
)

type GamePerson struct {
	name               [nameMaxSize]byte
	respectStrength    byte // [rrrrssss]
	levelExp           byte // [lllleeee]
	x, y, z, gold      int32
	manaHealth         [3]byte // [mm----hh][mmmmmmmm][hhhhhhhh]
	typeHouseGunFamily byte    // [xx------] type, [-----x--] house, [------x-] gun, [-------x] family
}

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i, c := range name {
			person.name[i] = byte(c)
		}

		if len(name) < nameMaxSize {
			person.name[len(name)] = byte(endOfString)
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealth[0] &= ^byte(manaFirstMask)
		person.manaHealth[0] |= byte((mana >> countOfBitsInByte) << manaFirstShift)
		person.manaHealth[1] = byte(mana & math.MaxUint8)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealth[0] &= ^byte(healthFirstMask)
		person.manaHealth[0] |= byte(health >> countOfBitsInByte)
		person.manaHealth[2] = byte(health & math.MaxUint8)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectStrength &= strengthMask
		person.respectStrength |= byte(respect) << respectShift
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectStrength &= respectMask
		person.respectStrength |= byte(strength & strengthMask)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levelExp &= levelMask
		person.levelExp |= byte(experience & experienceMask)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levelExp &= experienceMask
		person.levelExp |= byte(level) << levelShift
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= houseFlag
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= gunFlag
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= familyFlag
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily &= ^byte(typeMask)
		person.typeHouseGunFamily |= byte(personType) << typeShift
	}
}

func NewGamePerson(options ...Option) GamePerson {
	result := GamePerson{}

	for _, option := range options {
		option(&result)
	}

	return result
}

func (p *GamePerson) Name() string {
	length := 0
	for i := 0; i < nameMaxSize; i++ {
		if p.name[i] == endOfString {
			break
		}
		length++
	}
	return unsafe.String(&p.name[0], length)
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	low := int(p.manaHealth[1])
	top := int(p.manaHealth[0] >> manaFirstShift)
	return top<<countOfBitsInByte + low
}

func (p *GamePerson) Health() int {
	low := int(p.manaHealth[2])
	top := int(p.manaHealth[0] & healthFirstMask)
	return top<<countOfBitsInByte + low
}

func (p *GamePerson) Respect() int {
	return int(p.respectStrength >> respectShift)
}

func (p *GamePerson) Strength() int {
	return int(p.respectStrength & strengthMask)
}

func (p *GamePerson) Experience() int {
	return int(p.levelExp & experienceMask)
}

func (p *GamePerson) Level() int {
	return int(p.levelExp >> levelShift)
}

func (p *GamePerson) HasHouse() bool {
	return (p.typeHouseGunFamily & houseFlag) != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.typeHouseGunFamily & gunFlag) != 0
}

func (p *GamePerson) HasFamilty() bool {
	return (p.typeHouseGunFamily & familyFlag) != 0
}

func (p *GamePerson) Type() int {
	return int(p.typeHouseGunFamily >> typeShift)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)

	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
