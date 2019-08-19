package packinglist

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Input struct {
	Verbose      bool
	Days, People int
}

func PackingList(in *Input) (string, error) {
	var b strings.Builder
	for _, a := range activities {
		if a.Active() {
			a.Write(&b, 0)
		}
	}
	return b.String(), nil
}

var answers = map[string]bool{}

func query(s string) bool {
	// fmt.Println("query " + s)
	if a, ok := answers[s]; ok {
		return a
	}
	fmt.Print(s + " [y/N] ")
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		panic(err)
	}

	answers[s] = char == 'y'
	return answers[s]
}

type Activity struct {
	Name  string
	Check bool
	Items []*Activity
}

func (a *Activity) Active() bool {
	if !a.Check {
		return true
	}
	return query(a.Name)
}

func (a *Activity) Write(w io.Writer, l int) {
	pre := ""
	if l == 0 {
		pre = "\n"
	}
	pre += strings.Repeat("\t", l)
	w.Write([]byte(pre + "[ ] " + a.Name + "\n"))

	for _, i := range a.Items {
		if i.Active() {
			i.Write(w, l+1)
		}
	}
}

type Config func(a *Activity)

func act(name string, cfgs ...Config) *Activity {
	// fmt.Println("act " + name)
	a := &Activity{
		Name: name,
	}
	for _, cfg := range cfgs {
		cfg(a)
	}
	return a
}

func check(a *Activity) {
	a.Check = true
}

func subs(as ...*Activity) Config {
	// fmt.Println("subs")
	return func(a *Activity) {
		a.Items = append(a.Items, as...)
	}
}

func subItems(names ...string) Config {
	return func(a *Activity) {
		for _, n := range names {
			a.Items = append(a.Items, &Activity{Name: n})
		}
	}
}

func items(names ...string) []*Activity {
	var a []*Activity
	for _, n := range names {
		a = append(a, item(n))
	}
	return a
}

func item(n string) *Activity {
	return &Activity{Name: n}
}

var activities = []*Activity{
	{
		Name: "General",
		Items: items(
			"Notepad",
			"Pens",
			"Water Bottle",
			"Wallet",
			"Cash",
			"Day pack",
			"Sunglasses",
		),
	},
	{
		Name: "Electronics",
		Items: items(
			"Headphones",
			"Flashlight",
			"Phone Charger",
			"Watch Charger",
			"USB Multi Charger",
			"Book/Kindle",
		),
	},

	{
		Name: "Toiletries",
		Items: items(
			"Toothbrush",
			"Toothpaste",
			"Soap",
			"Shampoo",
			"Razer",
			"Shaving Cream",
			"Floss",
			"Deodorant",
			"Sunscreen",
		),
	}, {
		Name: "Misc",
		Items: items("Hand sanitizer",
			"Medicine",
			"First-aid kit",
			"Plastic baggies",
			"Laundry soap",
			"Clothesline",
			"Tissues",
			"Sewing kit",
			"Duct tape",
			"Insect repellant",
			"Packable duffel bag",
		),
	},

	{
		Name: "Clothing",
		Items: items(
			"Laundry Bag",
			"Pants/Shorts 2",
			"Shirts 5",
			"Jacket",
			"Underwear 5",
			"Socks 5",
			"Sweater",
			"Sleepwear",
			"Shoes/Boots",
		),
	},

	act("Swimming", check, subItems("swimsuit", "sandals", "rash/uv guard", "beach towel")),
	act("Formal event", check, subItems("Suit", "Tie", "Formal shoes", "Belt", "Collar stays")),

	act("Laptop", check, subItems("charger")),

	{
		Name:  "International",
		Check: true,
		Items: items("Passport", "Power adapter", "Pillow"),
	},
	{
		Name:  "Exercise",
		Items: items("rubber band", "trx", "lax ball", "stretch stick", "gym shorts", "gym shirts", "gym underwear", "gym shoes"),
	},
	act("Camping", check, subs(
		&Activity{
			Name:  "Camp",
			Items: items("chairs", "table", "garbage bags", "duct tape", "knife", "shovel", "toilet paper", "bleach", "solar blanket", "lantern", "hand sanitizer", "soap", "lighter", "matches", "batteries", "bug spray", "sunscreen", "rope", "tool roll"),
		},

		&Activity{Name: "Stove", Check: true, Items: items("propane", "connector")},

		&Activity{Name: "Coffee", Items: items("mugs", "grounds", "filter")},

		&Activity{
			Name:  "Sleeping",
			Items: items("Hammock", "Hammock Straps", "Tent", "Tarp", "Sleeping Bags", "pillows"),
		},

		&Activity{
			Name:  "Kitchen",
			Items: items("Cutting board", "knives", "sink", "dish soap", "pots and pans", "pot holders", "seasonings", "sponges", "serving spoon", "spatula", "tongs", "s'mores sticks", "scissors"),
		},

		&Activity{Name: "Eating", Items: items("water", "plates", "cups", "silverware", "bowls")},
	)),

	&Activity{
		Name:  "Work Travel",
		Check: true,
		Items: items("Badge", "Office Keys"),
	},
	{
		Name:  "Biking",
		Check: true,
		Items: items("Bike", "Helmet", "Gloves", "Pads", "Pump", "Shoes", "Tubes", "Water Bottle", "Lock", "Cable", "Keys"),
	},
	{
		Name:  "Snowboarding",
		Check: true,
		Items: items("Jacket", "Board", "Boots", "Pants", "Ninja Suits", "Helmet", "Goggles"),
	}, {
		Name:  "Fishing",
		Check: true,
		Items: items("Tackle Box", "Rod", "License", "Flies"),
	},
}
