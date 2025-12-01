package skill

import (
	"fmt"
)

const (
	Admin                 Skill = "Admin"
	Advocate              Skill = "Advocate"
	Animals_Handling      Skill = "Animals (Handling)"
	Animals_Veterinary    Skill = "Animals (Veterinary)"
	Animals_Training      Skill = "Animals (Training)"
	Art_Performing        Skill = "Art (Performing)"
	Art_Creative          Skill = "Art (Creative)"
	Art_Presentation      Skill = "Art (Presentation)"
	Astrogation           Skill = "Astrogation"
	Athletics_Dexterity   Skill = "Athletics (Dexterity)"
	Athletics_Endurance   Skill = "Athletics (Endurance)"
	Athletics_Strength    Skill = "Athletics (Strength)"
	Broker                Skill = "Broker"
	Carouse               Skill = "Carouse"
	Deception             Skill = "Deception"
	Diplomat              Skill = "Diplomat"
	Drive_Hovercraft      Skill = "Drive (Hovercraft)"
	Drive_Mole            Skill = "Drive (Mole)"
	Drive_Track           Skill = "Drive (Track)"
	Drive_Walker          Skill = "Drive (Walker)"
	Drive_Wheel           Skill = "Drive (Wheel)"
	Electronics_Comms     Skill = "Electronics (Comms)"
	Electronics_Computers Skill = "Electronics (Computers)"
	Electronics_RemoteOps Skill = "Electronics (Remote Ops)"
	Electronics_Sensors   Skill = "Electronics (Sensors)"
	Engineer_MDrive       Skill = "Engineer (M-drive)"
	Engineer_JDrive       Skill = "Engineer (J-drive)"
	Engineer_LiveSupport  Skill = "Engineer (Live Support)"
	Engineer_Power        Skill = "Engineer (Power)"
	Explosives            Skill = "Explosives"
	Flyer_Airship         Skill = "Flyer (Airship)"
	Flyer_Grav            Skill = "Flyer (Grav)"
	Flyer_Ornithopter     Skill = "Flyer (Ornithopter)"
	Flyer_Rotor           Skill = "Flyer (Rotor)"
	Flyer_Wing            Skill = "Flyer (Wing)"
	Gambler               Skill = "Gambler"
	Gunner_Turret         Skill = "Gunner (Turret)"
	Gunner_Ortilery       Skill = "Gunner (Ortilery)"
	Gunner_Screen         Skill = "Gunner (Screen)"
	Gunner_Capital        Skill = "Gunner (Capital)"
	GunCombat_Archaic     Skill = "Gun Combat (Archaic)"
	GunCombat_Energy      Skill = "Gun Combat (Energy)"
	GunCombat_Slug        Skill = "Gun Combat (Slug)"
	HeavyWeapons_Artilery Skill = "Heavy Weapons (Artilery)"
	HeavyWeapons_Portable Skill = "Heavy Weapons (Portable)"
	HeavyWeapons_Vechicle Skill = "Heavy Weapons (Vechicle)"
	Independence          Skill = "Independence"
	Investigate           Skill = "Investigate"
	JOAT                  Skill = "Jack-of-All-Trades"
	Language              Skill = "Language"
	Leadership            Skill = "Leadership"
	Mechanic              Skill = "Mechanic"
	Medic                 Skill = "Medic"
	Melee_Unarmed         Skill = "Melee (Unarmed)"
	Melee_Blade           Skill = "Melee (Blade)"
	Melee_Bludgen         Skill = "Melee (Bludgen)"
	Melee_Natural         Skill = "Melee (Natural)"
	Navigation            Skill = "Navigation"
	Persuade              Skill = "Persuade"
	Pilot_SmallCraft      Skill = "Pilot (Small Craft)"
	Pilot_Spaceship       Skill = "Pilot (Spaceship)"
	Pilot_CapitalShips    Skill = "Pilot (Capital Ships)"
	Profession            Skill = "Profession"
	Recon                 Skill = "Recon"
	ScienceLife           Skill = "Science (Life)"
	SciencePhysical       Skill = "Science (Physical)"
	ScienceRobotics       Skill = "Science (Robotics)"
	ScienceSocial         Skill = "Science (Social)"
	ScienceSpaces         Skill = "Science (Space)"
	Seafarer              Skill = "Seafarer"
	Stealth               Skill = "Stealth"
	Steward               Skill = "Steward"
	Streetwise            Skill = "Streetwise"
	Survival              Skill = "Survival"
	Tactics_Military      Skill = "Tactics (Military)"
	Tactics_Navy          Skill = "Tactics (Navy)"
	Tactics_SmallUnit     Skill = "Tactics (Small Unit)"
	Tolerance             Skill = "Tolerance"
	VaccSuit              Skill = "Vacc Suit"
)

type Skill string

func (s Skill) String() string {
	return fmt.Sprintf("%s", s)
}

func BackgroundSkillList() []Skill {
	return []Skill{
		Admin,
		Animals_Handling,
		Animals_Training,
		Animals_Veterinary,
		Art_Creative,
		Art_Performing,
		Art_Presentation,
		Athletics_Dexterity,
		Athletics_Endurance,
		Athletics_Strength,
		Carouse,
		Drive_Hovercraft,
		Drive_Mole,
		Drive_Track,
		Drive_Walker,
		Drive_Wheel,
		Electronics_Comms,
		Electronics_Computers,
		Electronics_RemoteOps,
		Electronics_Sensors,
		Flyer_Airship,
		Flyer_Grav,
		Flyer_Ornithopter,
		Flyer_Rotor,
		Flyer_Wing,
		Language,
		Mechanic,
		Medic,
		Profession,
		ScienceLife,
		SciencePhysical,
		ScienceRobotics,
		ScienceSocial,
		ScienceSpaces,
		Streetwise,
		Survival,
		VaccSuit,
	}
}
