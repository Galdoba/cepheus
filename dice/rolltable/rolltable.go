package rolltable

import (
	"fmt"

	"github.com/Galdoba/t5/pkg/dice"
	"github.com/Galdoba/t5/pkg/dice/index"
)

type Table struct {
	Description string
	DiceCode    string
	Entries     map[string]TableEntry
}

func NewTable(options ...TableOption) Table {
	t := Table{}
	for _, modify := range options {
		modify(&t)
	}
	return t
}

type TableOption func(*Table)

func WithEntries(entries ...TableEntry) TableOption {
	return func(t *Table) {
		if len(t.Entries) < 1 {
			t.Entries = make(map[string]TableEntry)
		}
		for _, entry := range entries {
			_, err := index.Parse(entry.Key)
			if err != nil {
				panic(err.Error())
			}
			t.Entries[entry.Key] = entry
		}
	}
}

func WithDescription(desc string) TableOption {
	return func(t *Table) {
		t.Description = desc
	}
}

func WithDiceCode(code string) TableOption {
	return func(t *Table) {
		t.DiceCode = code
	}
}

func (t Table) Check(dp *dice.Dicepool, opts ...dice.RollOption) TableEntry {
	i := dp.Sum(t.DiceCode, opts...)
	fmt.Println("roll", i)
	for _, entry := range t.Entries {
		for _, index := range entry.Indexes {
			if i == index {
				return entry
			}
		}
	}
	return NoEntry
}

type TableEntry struct {
	Children          []Table
	Key               string
	Indexes           []int
	EntryText         string
	ResultValueString map[string]string
	ResultValueInt    map[string]int
	ResultValueBool   map[string]bool
	ResultValueBytes  map[string][]byte
}

var NoEntry = TableEntry{
	EntryText: "No Entry",
}

func NewEntry(text string, options ...EntryOption) TableEntry {
	te := TableEntry{}
	te.EntryText = text
	for _, modify := range options {
		modify(&te)
	}
	return te
}

type EntryOption func(*TableEntry)

func WithIndexes(inds ...int) EntryOption {
	return func(te *TableEntry) {
		key, err := index.Format(inds)
		if err != nil {
			panic(err)
		}
		te.Key = key
		te.Indexes = append(te.Indexes, inds...)
	}
}

func WithKey(key string) EntryOption {
	return func(te *TableEntry) {
		inds, err := index.Parse(key)
		if err != nil {
			panic(err)
		}
		te.Key = key
		te.Indexes = append(te.Indexes, inds...)
	}
}

func WithStringValue(key, val string) EntryOption {
	return func(te *TableEntry) {
		if te.ResultValueString == nil {
			te.ResultValueString = make(map[string]string)
		}
		te.ResultValueString[key] = val
	}
}

func WithIntValue(key string, val int) EntryOption {
	return func(te *TableEntry) {
		if te.ResultValueString == nil {
			te.ResultValueInt = make(map[string]int)
		}
		te.ResultValueInt[key] = val
	}
}

func WithBoolValue(key string, val bool) EntryOption {
	return func(te *TableEntry) {
		if te.ResultValueString == nil {
			te.ResultValueBool = make(map[string]bool)
		}
		te.ResultValueBool[key] = val
	}
}

func WithBytesValue(key string, val []byte) EntryOption {
	return func(te *TableEntry) {
		if te.ResultValueString == nil {
			te.ResultValueBytes = make(map[string][]byte)
		}
		te.ResultValueBytes[key] = val
	}
}

func WithChildTable(tables []Table) EntryOption {
	return func(te *TableEntry) {
		te.Children = tables
	}
}
