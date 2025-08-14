package action

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/interaction"
	"github.com/urfave/cli/v3"
)

const (
	phase_SETUP = iota
	phase_RESOLUTION
	phase_CONCLUSION
	status_CONTINUE
	status_DEFENDER_WINS
	status_ATTACKER_WINS
	status_SHIP_DESTROYED
)

type Boarding struct {
	PowerBalance    int //modifier to main roll: positive good for attacker, negative - for defender
	PowerAdjustment int //modifier to main roll: positive good for attacker, negative - for defender
	ShipMaxHull     int
	ShipCurrentHull int
	ShipArmor       int
	BoardingRound   int
	CombatRound     int
	MinutesPassed   int
	dp              *dice.Dicepool
	phase           int
	interactive     bool
	aggressionLevel int
	status          int
	log             []string
}

func BoardingAction(ctx context.Context, c *cli.Command) error {
	brd, err := SetupBoardingAction(c)
	if err != nil {
		return err
	}
	brd.Log("")
	brd.Log("RESOLUTION PHASE:")
	for brd.status == status_CONTINUE {
		brd.Log("Boarding Act: %v; start Combat Round: %v; %v passed", brd.BoardingRound, brd.CombatRound, timeStr(brd.MinutesPassed))
		playersEffect := getNumber("Player's Actions Effect:", "Effect of player's task chain")
		dr := brd.dp.Sum("2d6")
		ar := brd.dp.Sum("2d6")
		brd.Log("attacker roll: %v; defender roll: %v", ar, dr)
		effect := ar - dr + brd.PowerBalance + playersEffect
		brd.Log("Act Resolution:")
		brd.Log("effect: %v = (%v - %v + (%v) + (%v) + (%v))", effect, ar, dr, brd.PowerBalance, brd.PowerAdjustment, playersEffect)
		actRounds, actMinutes := brd.addRounds(2)
		brd.Log("action took: %v Combat Rounds (%v)", actRounds, timeStr(actMinutes))
		hullDecrease := 0
		brd.PowerAdjustment = 0
		switch bound(effect, -7, 7) {
		case -7:
			brd.status = status_DEFENDER_WINS
			brd.Log("The attackers are soundly defeated. If the attacker’s ship is docked with the defender’s, the defenders may mount a new boarding action of their own and gain DM+4 on the roll to resolve it.")
		case -4, -5, -6:
			brd.status = status_DEFENDER_WINS
			brd.Log("The attackers are defeated. The attackers must retreat back to their own ship or space – if they are unable to do so, they are killed or captured.")
		case -1, -2, -3:
			actRounds, _ := brd.addRounds(1)
			brd.PowerAdjustment += -2
			hullDecrease = brd.dp.Sum("2d6")
			brd.Log("Report: Fighting continues. Resolve the boarding again in another %v rounds but defender gains +2 DM to their roll.", actRounds)
			brd.applyDamage(hullDecrease, true)
		case 0:
			actRounds, _ := brd.addRounds(1)
			// brd.PowerAdjustment = 0
			brd.Log("Report: Fighting continues. Resolve the boarding again in another %v rounds.", actRounds)
		case 1, 2, 3:
			actRounds, _ := brd.addRounds(1)
			brd.PowerAdjustment += 2
			hullDecrease = brd.dp.Sum("2d6")
			brd.Log("Report: Fighting continues. Resolve the boarding again in another %v rounds but attacker gains +2 DM to their roll.", actRounds)
			brd.applyDamage(hullDecrease, true)
		case 4, 5, 6:
			brd.status = status_ATTACKER_WINS
			hullDecrease = brd.dp.Sum("1d6")
			actRounds, _ := brd.addRounds(2)
			brd.Log("The boarding action is successful and the ship being boarded suffers %v damage, ignoring any armour. The attackers may take control of the ship after another %v rounds of pacification.", hullDecrease, actRounds)
		case 7:
			brd.status = status_ATTACKER_WINS
			brd.Log("The attackers storm the enemy ship and take control of it immediately.")
		default:
			return fmt.Errorf("unexpected effect result: %v", bound(effect, -7, 7))
		}
		if brd.ShipCurrentHull < 1 {
			brd.status = status_SHIP_DESTROYED
		}
		brd.BoardingRound++
		brd.Log("")
	}
	brd.phase = phase_CONCLUSION
	brd.Log("After a battle that lasted %v (%v Combat Rounds)", timeStr(brd.MinutesPassed), brd.CombatRound-1)
	switch brd.status {
	case status_ATTACKER_WINS:
		brd.Log("Attackers took the ship with %v hull points left.", brd.ShipCurrentHull)
	case status_DEFENDER_WINS:
		brd.Log("Defender kept the ship with %v hull points left.", brd.ShipCurrentHull)
	case status_SHIP_DESTROYED:
		brd.Log("Ship became comletly inoperable")

	}
	return nil

}

func SetupBoardingAction(c *cli.Command) (*Boarding, error) {
	brd := Boarding{}
	brd.Log("SETUP PHASE:")
	brd.PowerBalance = getNumber("Set Power Balance:", "positive is good for attackers, negative is good for defenders")
	brd.Log("Power balanse is set to %d", brd.PowerBalance)
	brd.ShipMaxHull = getNumber("Set Defending Ship Maximum Hull points:", "")
	brd.ShipCurrentHull = getNumber("Set Defending Ship Current Hull points:", "")
	if brd.ShipCurrentHull > brd.ShipMaxHull {
		brd.ShipCurrentHull, brd.ShipMaxHull = brd.ShipMaxHull, brd.ShipCurrentHull
	}
	brd.ShipArmor = getNumber("Set Ship Armor:", "")
	msg := fmt.Sprintf("Boarded ship condition: %v of %v Hull points", brd.ShipCurrentHull, brd.ShipMaxHull)
	brd.Log(msg)
	brd.dp = dice.NewDicepool()
	brd.phase = phase_RESOLUTION
	brd.status = status_CONTINUE
	brd.aggressionLevel = brd.dp.Sum("2d6")
	brd.BoardingRound = 1
	brd.CombatRound = 1
	return &brd, nil
}

func (brd *Boarding) Log(format string, args ...any) {
	format += "\n"
	fmt.Fprintf(os.Stderr, format, args...)
	brd.log = append(brd.log, fmt.Sprintf(format, args...))
}

func getNumber(title, descr string) int {
	input, err := interaction.GetInput(title,
		interaction.WithValidator(interaction.Number),
		interaction.WithDescription(descr),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	n, _ := strconv.Atoi(input)
	return n
}

func persentage(top, down int) int {
	return top * 100 / down
}

func timeStr(minutes int) string {
	h := 0
	for minutes > 59 {
		h += 1
		minutes -= 60
	}
	switch h {
	case 0:
		return fmt.Sprintf("%d minutes", minutes)
	case 1:
		return fmt.Sprintf("1 hour %d minutes", minutes)
	default:
		return fmt.Sprintf("%d hours %d minutes", h, minutes)
	}

}

func bound(i, min, max int) int {
	if i < min {
		return min

	}
	if i > max {
		return max
	}
	return i
}

func (brd *Boarding) addRounds(dices int) (int, int) {
	actRounds := brd.dp.Sum(fmt.Sprintf("%vd6", dices))
	actMinutes := actRounds * 6
	brd.CombatRound += actRounds
	brd.MinutesPassed += actMinutes
	return actRounds, actMinutes
}

func (brd *Boarding) applyDamage(damage int, useArmor bool) {
	damageReceived := 0
	switch useArmor {
	case true:
		consumed := 0
		switch damage < brd.ShipArmor {
		case true:
			consumed = damage
		default:
			consumed = brd.ShipArmor
		}
		damageReceived = bound(damage-consumed, 0, damage)
		brd.Log("Ship received %v Hull damage (%v was consumed by Armor)", damageReceived, brd.ShipArmor)
	case false:
		damageReceived = damage
		brd.Log("Ship received %v Hull damage, ignoring any armor.", damageReceived)
	}
	brd.ShipCurrentHull = bound(brd.ShipCurrentHull-damageReceived, 0, brd.ShipCurrentHull)
}
