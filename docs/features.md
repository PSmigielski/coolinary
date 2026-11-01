# Coolinary — Feature Plan

Scope: what the app does, split by whether a feature needs AI. No tech-stack choices. Builds on `docs/business-model.md`.

## Context

The business model sets two cost limits: a free user may cost at most €0.10 a month, and AI may take at most 15% of what a paying household brings in. So every feature is placed in one of three groups:

| Group | Meaning | Cost | Who gets it |
|---|---|---|---|
| **Without AI** | Rules and stored data | None per use; works offline and when AI is down | Everyone |
| **AI, run once for everyone** | AI prepares shared content centrally (weekly flyers, recipe library, translations) | Fixed, does not grow with users | Everyone |
| **AI, run per request** | AI works on one user's photo or fridge | Grows with use | Capped on Free, unlimited on Plus |

Two rules follow:

1. **The app stays useful with AI switched off.** A user can still keep a fridge list, get dish matches from the recipe library, filter by cuisine and tools, and build a shopping list from flyers.
2. **AI proposes, rules check.** Every AI-composed dish is checked without AI against the user's tools, allergens, and stock before it is shown. Calories and portion sizes are always calculated without AI, from nutrition tables; AI never states a calorie number.

## 1. Kitchen profile (all without AI)

Set once at first start, with presets so it takes four screens at most. People and goals are optional and are not asked at first start.

| Setting | What the user picks | Effect |
|---|---|---|
| **Cuisines** | Favourites from: Polish, Italian, Asian (optionally narrowed to Chinese, Japanese, Thai, Indian, Vietnamese, Korean), Mediterranean, French, Spanish, German, Mexican, Middle Eastern, American, or "no preference". Can differ per meal (Polish breakfast, Asian dinner). Home cuisine is pre-selected by country. | Filters the library and steers AI composition |
| **Kitchen tools** | Checklist: hob, oven, microwave, air fryer, blender, hand blender, food processor, multicooker / pressure cooker, cooking robot, toaster / sandwich maker, grill, steamer, rice cooker, mixer, wok, baking tins. Presets: "student room" (microwave + kettle), "basic kitchen", "fully equipped". | Hard limit: no dish ever needs an unticked tool |
| Diet and allergens | Vegetarian, vegan, gluten-free, lactose-free, nut-free, dislikes | Hard limit |
| Household | Portions, which meal is the main hot one, time per meal, and how often to cook: every day, every other day, or twice a week | Shapes the day and week plan |
| People and goals | Optional, for each person who eats at home, children included (a child needs no account): age, and a goal: "just show calories", "keep my weight", "lose slowly", "gain", or "my own number" from a doctor or dietitian. Children's goals are set by a parent and are different (see "Planning for each person"). Height, weight, and activity only when the person wants an estimate. | Sets each person's daily calorie target and portion size (see "Planning for each person" in section 2) |
| Staples | "Always at home" items (salt, oil, flour) | Never counted as missing |
| My shops and location | Shops suggested from the user's location (or a typed postcode), nearest first; the user confirms or edits them. Travel range: walking (1 km), short trip (3 km), by car (10 km). "How you shop": how far apart two shops may be (none up to 1.5 km), how many shops one trip may have (1 to 3), and how much an extra shop must save (default €2). | Selects which flyers are used, how offers are ordered, and how the shopping trip is planned |

## 2. Features by area

### Fridge and pantry
**Without AI**
- Add items by typing, with autocomplete from an ingredient catalogue.
- Barcode scan: the code fills in the product's name and size, and the app suggests the storage place and use-by date (see section 4).
- Unpack mode: scan a whole bag of barcodes in a row and confirm the list once.
- Quantity, storage place (fridge, freezer, pantry), expiry date with a default shelf life per product type.
- Expiry reminders.
- Stock is reduced automatically when a dish is marked "cooked".

**With AI**
- Photo of the fridge or a shelf becomes an item list the user confirms (per request).
- Receipt photo becomes items (per request).
- A typed or spoken sentence ("bought eggs, two peppers, some cheese") becomes items (per request).
- Ingredient catalogue translated and matched across languages (once for everyone).
- Each newly seen product is matched to a catalogue ingredient, so recipes can use it (once for everyone).

### Dishes and the day plan
**Without AI**
- Recipe library; each dish tagged by meal (breakfast, dinner, supper), cuisine, tools, time, and diet.
- "What can I make now": library dishes ranked by share of ingredients on hand and by what expires first.
- Filters, favourites, history, ratings, portion scaling.
- Step-by-step cooking mode with timers and a "cooked it" button.

**With AI**
- Compose a full day (breakfast, dinner, supper) from current stock, oldest items first, so the three meals do not compete for the same ingredients (per request).
- Swap one meal, replace a missing ingredient, or follow a free-text wish such as "something warm in 20 minutes" (per request).
- Week plan, including dinner leftovers reused in the next supper (per request).
- Well-rated AI dishes are reviewed and added to the library, so the same dish later costs nothing (once for everyone).

### Calories and goals
**Without AI**
- Every dish shows calories per portion ("about 520 kcal"). One tap shows protein, carbohydrate, sugars, fat, fibre, and salt. The values come from nutrition tables (section 5) and change with portion scaling.
- kcal by default; kJ as an option in settings.
- When goals are set, each dish and cooking mode show every person's portion, and the day plan shows each person's day total next to their target.
- Each person can hide calories. On children's profiles they are hidden by default and only portions are shown.
- Library filters state numbers: calories per portion (under 400, 400–600, over 600) and "25 g protein or more". No dish is labelled "healthy" or "light".
- Effort tags on library dishes: no cooking, keeps 3 days in the fridge, freezes well, good cold the next day, base plus toppings.
- A person can be marked away for a meal (school lunch, eating out). The pot shrinks and the shopping list follows.
- Coolinary plans meals to the targets. It is not a food diary: food eaten outside the plan is not logged.

**With AI**
- Compose a family day with one dish per meal that fits every person's target, from stock first (per request).
- Week plan with the fewest cooking sessions: cook a double amount on cooking days and plan the rest for a later meal (per request).
- Match each catalogue ingredient to its entry in the nutrition tables, and suggest effort tags for library dishes; both are checked by hand before use (once for everyone).

AI receives only each person's daily target, never their weight. The portion arithmetic runs without AI, so a family plan costs about the same as a plan for one.

### Planning for each person (without AI)

The aim is that each person eats the right amount while the kitchen cooks as little as possible. The family shares one dish per meal and each person gets a portion of a different size. Nobody gets a separate dish because of calories.

**Each person's daily target**
- Adults: resting energy need from age, sex, height, and weight (Mifflin–St Jeor equation), multiplied by an activity level: mostly sitting 1.4, some activity 1.6, active 1.8.
- "Lose slowly" takes 500 kcal a day off the estimate; "gain" adds 300.
- No adult target is ever below 1,200 kcal a day, including "my own number". A lower number is refused, and the person is asked to follow their doctor's plan.
- The estimate is shown as "a starting point". Weight can be updated to refresh it, but no weight history is kept.

**Children and teenagers under 18**

About 1 in 4 children aged 7–9 in Europe live with overweight and 11% with obesity, so the app has to help these families. Guidance from NICE and the AAP does not recommend a calorie cut for a child. It recommends changing what the whole family eats, with the parent responsible, and for many children it aims to keep weight steady while they grow taller. A child's needs also change as they grow, and adult formulas do not fit them. So a child gets no "lose" goal calculated by a formula. A parent picks one of these:

| Goal (set by a parent) | Target | When |
|---|---|---|
| Growing (default) | EFSA reference value for the child's age, sex, and activity | Most children |
| Grow into a healthier weight | EFSA reference value for the child's age and sex at the lowest activity level, never lower. EFSA's values assume a healthy weight, so a heavier child eats like a healthy child of the same age and grows into their weight | The parent is worried about the child's weight |
| Set by our doctor or dietitian | The number from the child's doctor or dietitian. Numbers below a lower limit for the child's age are refused | The child is in a weight programme |

- **The whole family eats lighter together.** When any child has the second or third goal, the app offers "eat lighter as a family". The planner then prefers dishes with more vegetables and fewer calories per portion, and uses fruit, vegetables, milk, and yoghurt as no-cook extras for everyone instead of bread, cheese, and nuts. Each adult keeps their own target. This matches the guidance that every family member eat healthily, whatever their weight.
- **The child is never singled out.** The child sees only their portion: no goal, no calories, no "diet" wording. Everyone eats the same dish.
- **Only a parent sets the goal.** A teenager with their own account cannot change it.
- **No weight history, no before-and-after,** for children as for adults. The 2019 backlash against Kurbo, a weight-loss app for 8 to 17 year olds, centred on food tracking and children's before-and-after stories.
- When "set by our doctor or dietitian" is chosen, the app suggests sharing the week plan with that doctor or dietitian.

**How one meal is split**
1. Each person's day target is shared over the meals: 25% breakfast, 40% the main hot meal, 25% the other meal, and 10% left for snacks and drinks. The household can change the shares.
2. People marked away are left out of that meal.
3. Each person's portion is their share divided by the calories of one standard portion, rounded to the nearest quarter, between ½ and 2 portions.
4. The pot is the sum of everyone's portions, rounded up to a whole portion. Ingredients and the shopping list are scaled to it.
5. If the 2-portion limit leaves someone more than 15% short, that person gets a no-cook extra: bread, fruit, yoghurt, milk, cheese, nuts, or raw vegetables, taken from what is at home first. If even ½ portion leaves someone more than 15% over, the planner picks a different dish.
6. Smaller gaps carry over: a bigger dinner makes the next meal smaller. Each person's day should land within 10% of their target, and their week within 5%.
7. The dish card and cooking mode say how to divide the pot: "6 portions: Adam 1½, Ewa 1½, Ola 1, Kuba 2".

**One dish per meal, few cooking sessions**
- The shared dish meets the allergens and diet of everyone eating it. When one person's diet narrows the choice too much (one vegetarian, for example), the planner prefers "base plus toppings" dishes: one base for all and toppings chosen per plate (tacos, bowls, wraps). A topping that needs its own pan counts as extra cooking.
- Breakfast and the meal that is not the main one come from "no cooking" dishes unless the household asks otherwise.
- On a cooking day the main dish is cooked in a double amount if it keeps 3 days or freezes well. The second half becomes a later meal and is split again for whoever eats then.
- The week plan shows the effort: "3 cooking sessions this week, about 2 hours at the stove".
- When these aims conflict, this order wins: hard limits (allergens, diet, tools); then the cooking rhythm the household chose; then each person's day within 10%; then oldest stock first; then this week's offers.

**Privacy**
- Age, height, weight, and goals count as health data under the GDPR. They are asked for with a separate, explicit consent and can be deleted at any time.
- They stay on the phone, and the target is calculated there. Only each person's daily calorie target is stored on the server.
- Other household members see portions only, never another adult's weight or goal. Parents see their children's.

**Worked example** (invented targets; dinner is the main meal; chicken curry with rice at 520 kcal per standard portion)

| Person | Goal | Day target | Dinner share (40%) | Portion | kcal | Extra |
|---|---|---|---|---|---|---|
| Adam | Lose slowly | 2,000 | 800 | 1½ | 780 | none |
| Ewa | Keep | 1,800 | 720 | 1½ | 780 | none; supper is 60 kcal smaller |
| Ola, 8 | Growing | 1,400 | 560 | 1 | 520 | none |
| Kuba, 15, plays football | Growing, active | 3,200 | 1,280 | 2 (the most) | 1,040 | a glass of milk and a banana, about 260 kcal |

The pot is 6 portions, cooked once. Kuba is 240 kcal (19%) short at 2 portions, so only he gets an extra. With "every other day" cooking the pot becomes 12 portions, and the second half is tomorrow's dinner, split again the same way.

### Cuisine choice
**Without AI**
- Cuisine works as a filter on the library.
- Each cuisine has a "base pantry" checklist (Asian: soy sauce, rice, ginger, sesame oil; Italian: olive oil, passata, pasta, hard cheese). One tap adds the missing base items to the shopping list.

**With AI**
- Compose in the chosen cuisine from what is in the fridge (per request).
- When stock does not fit the cuisine, offer three choices: the closest dish possible now, the 1–3 items to buy, or a different cuisine for today (per request).
- Restyle a dish ("make this Italian-style") (per request).

### Kitchen tools
**Without AI**
- The checklist from the profile filters every library dish.
- One-off filter for tonight ("air fryer only").

**With AI**
- Rewrite a recipe for the tools at hand: oven dish as a pan version, air-fryer times, cooking-robot steps (per request).

### Flyers and offers
**Without AI**
- **Nearest shops first:** offers are ordered by distance from the user, and only shops inside the chosen travel range are shown. Each offer shows the shop and its distance.
- Location also picks the right edition where a chain's flyer differs by store or region (Kaufland, Aldi Nord or Süd).
- Browse and search this week's offers by shop, category, and end date.
- Tap an offer to add it to the shopping list.
- "Offer ends tomorrow" reminders.
- Every offer links to the original flyer.

**With AI**
- Each week, read every covered chain's flyer once and extract the offers: product, size, price, conditions ("2+1", "with app or card"), valid dates; match product names to catalogue ingredients (once for everyone).
- "Plan from deals": compose dishes from the fridge plus this week's offers, preferring the nearest shops (per request).
- Photograph or share any flyer that is not covered yet and get its offers extracted (per request). This works in every country from the first flyer release.

### Shopping list
**Without AI: how the list is made**
1. Take every ingredient of the dishes in the chosen day or week plan, scaled to the household's portions, or to the pot size when goals are set (see "Planning for each person").
2. Subtract what is in the fridge, freezer, and pantry, and the "always at home" staples.
3. Merge the same ingredient from several dishes into one line, and round it up to a whole pack where the catalogue knows the pack size.
4. Add the extras: base items for a chosen cuisine, and anything a household member adds by hand.
5. Match each line to this week's offers at shops in range and show the offer badge ("−30% at Lidl until Saturday").
6. Plan the trip: choose the shop or shops to buy at (see "Planning the shopping trip" below), then order the list by that route, grouped by shop section, with a subtotal per shop.
7. When an item is ticked in the shop, it moves into the fridge with a default use-by date.

**Also without AI**
- The user reviews the proposed lines before they join the list and can untick any of them.
- Each line shows where it came from: which meal, which cuisine base, or who added it.
- The list is shared live between household members.
- Optional, later: a reminder when the user is near a shop that has list items on offer.

**With AI**
- "Cheapest top-up": the smallest basket that completes the week plan using current offers, including swapping an ingredient for a discounted equivalent (per request).

### Planning the shopping trip (without AI)

The list is turned into a trip by a small optimisation from the same family as the knapsack problem: find the shops that give the lowest bill without exceeding the walking the user accepts. The aim is to buy as much as possible in one shop and add another only when it pays.

**What the user sets** (the "How you shop" screen)
- Range from home: 1, 3, or 10 km.
- Walk between shops: from none (one shop only) up to 1.5 km.
- Shops per trip: 1 to 3.
- Minimum saving for an extra shop: €2 by default.

**How the plan is found**
1. Take the shops in range that the app has prices for.
2. Build every set of shops, from one up to the allowed number, in which each walk from one shop to the next stays within the limit. With 8 shops and at most 3 per trip that is 92 sets, so all of them are checked.
3. In each set, give every list line to the shop where it costs least, and add up the bill.
4. Start from the cheapest single shop. A larger set replaces it only if it saves at least the minimum for each extra shop.
5. Ties go to fewer shops, then to the shorter walk.
6. Show the winning plan and up to two alternatives with their price difference and extra distance. The user can pick any of them.

**Where prices come from**
- The flyer price, when the item is on offer at that shop this week.
- Otherwise the last price seen on receipts that users scanned at that chain (from release 3).
- Otherwise a reference price that is the same for every shop, so it does not sway the choice.
- Totals are labelled "about". Until receipt prices build up, only flyer offers separate one shop from another.

**Distance** is the straight line between shop positions multiplied by 1.3 as a walking estimate, calculated on the phone.

**Worked example** (invented prices in euros; Lidl is 350 m from home, Kaufland 1.2 km, Aldi 2.4 km; Lidl to Kaufland is 900 m)

| Item | Lidl | Kaufland | Aldi |
|---|---|---|---|
| Chicken thighs 1 kg | 3.89 | 5.49 | 5.29 |
| Basmati rice 1 kg | 1.79 | 2.29 | 2.19 |
| Onions 1 kg | 1.49 | 0.99 | 1.39 |
| Soy sauce 150 ml | 2.49 | 1.59 | 2.29 |
| Ginger | 0.99 | 0.89 | 0.95 |
| **Whole list in this shop** | **10.65** | **11.25** | **12.11** |

Lidl then Kaufland costs 9.15, which saves 1.50 over Lidl alone. With the default €2 minimum the plan is Lidl only; with a €1 minimum it is Lidl then Kaufland; with the walk limit under 900 m the two-shop plan is not offered at all.

### Savings and household (all without AI)
- Counters for euros and kilograms saved (used before expiry, bought on offer) and a monthly report.
- Notifications, language, units.

### Linking people into a household (all without AI)
- A household is one shared fridge, one shopping list, and one meal plan for up to 5 people.
- Children, and anyone without a phone, can be added as profiles without an account, for portions and goals only. These profiles do not count towards the 5.
- An invite can be handed over three ways: a QR code shown on screen, a 6-character code, or a link sent through any messaging app.
- An invite works for 48 hours and for one person. Waiting invites are listed and can be cancelled.
- The person joining sees whose household it is before confirming, and chooses whether to bring their own fridge items along.
- The owner can remove people. Anyone can leave, and then starts again with their own empty fridge.
- Sharing needs Plus, held by the owner. On Free, "Invite someone" opens the Plus offer.

## 3. Where flyer data comes from

Used in this order:

1. **Flyers the user supplies** (photo, PDF, or link). Any country, lowest legal risk.
2. **Central reading of each chain's own flyer page** (sources below). Only facts are stored (product, price, dates, shop) plus a link to the source; no flyer pages are republished.
3. **Licensed feed from an offers network** such as ShopFully (formerly Offerista; 25+ countries, 1,400+ partner apps), once there is traffic to show. This may also bring revenue.

Each retailer's terms must be checked before step 2 goes live for that chain. EU case law (Ryanair v PR Aviation) lets a site restrict reuse of its data through its terms, so this is a real gate. Flyer apps (Blix, Moja Gazetka, kaufDA, marktguru) are not used as sources: their collections are their own product.

### Sources per shop: launch chains (release 2)

All checked on 4 October 2026.

| Shop | Flyer list page | Individual flyers | Edition | Rhythm seen |
|---|---|---|---|---|
| Biedronka (PL) | `biedronka.pl/pl/gazetki` | `biedronka.pl/pl/press,id,[id],title,[name]`; page images on `cdn.biedronka.pl` | National | Several flyers a week; start date is in the title |
| Lidl (PL) | `lidl.pl/c/nasze-gazetki/s10008614` | Online flipbook | National | Monday and Thursday flyers |
| Kaufland (PL) | `sklep.kaufland.pl/gazeta-reklamowa.html` | `leaflets.kaufland.com/pl-PL/[id]/ar/[store]` | Per store version | Thursday to Wednesday; online from Tuesday |
| Lidl (DE) | `lidl.de/c/online-prospekte/s10005610` | `lidl.de/l/prospekte/[name]/ar/0` | National | Monday to Saturday; three weeks published ahead |
| Aldi Süd (DE) | `aldi-sued.de/de/angebote/prospekte.html` | `prospekt.aldi-sued.de/kw[week]-[year]-…` | Southern and western Germany | This week plus two weeks ahead |
| Aldi Nord (DE) | `aldi-nord.de/prospekte.html` | `magazine.aldi-nord.de/aldi-nord/aldi-aktuell/[year]/…`; PDF offered | Northern and eastern Germany, plus regional extras | Monday to Saturday; up to 14 days ahead |
| Kaufland (DE) | `filiale.kaufland.de/prospekte.html` | `leaflets.kaufland.com/de-DE/[id]/ar/[store]` | Per store | Thursday to Wednesday |

Kaufland's flyer depends on the store and Aldi is two separate companies split by region, so the user's location decides which edition they see from release 2 onward.

### Sources per shop: next chains (release 3)

| Shop | Flyer list page | Notes |
|---|---|---|
| Dino (PL) | `marketdino.pl/gazetki/lista/standard` | Weekly, Wednesday to Tuesday |
| Netto (PL) | `netto.pl/gazetka-netto/` | Flyers at `?leafletid=[id]&page=[n]`; Monday and Thursday |
| Carrefour (PL) | `carrefour.pl/gazetka-handlowa` | PDF offered; differs by store format |
| Auchan (PL) | `auchan.pl/pl/gazetki` | National plus per-store flyers |
| Aldi (PL) | `aldi.pl/gazetki.html` | Page exists; details not yet checked |
| Penny (DE) | `penny.de/angebote` | Monday to Saturday; next week shown from Friday; extra local offers per store |
| Rewe, Edeka, Netto (DE) | `rewe.de/angebote`, `edeka.de`, `netto-online.de/filialangebote` | Per-store offers, and all three sites refused automated access in the check. Covered only through the licensed feed or user-supplied flyers. |

### Keeping flyers in sync

1. **Schedule:** each list page is checked once a day, plus an extra check just after the chain's usual publishing time (for example Tuesday and Wednesday for Kaufland).
2. **Detect:** the flyer IDs and dates on the list page are compared with what is stored. A new ID means a new flyer; changed dates or page count mean a re-read.
3. **Fetch:** each flyer's pages are downloaded once, slowly, with the app identified in the request and the site's robot rules respected.
4. **Read:** AI extracts the offers (once for everyone).
5. **Check without AI:** price above zero, dates inside the flyer's validity, unit present, product matched to the catalogue; anything that fails is held back. The founder spot-checks 20 offers per chain each week.
6. **Publish:** offers go live on their start date. Flyers published ahead appear as "next week", so users can plan before the offers start.
7. **Expire:** offers disappear on their end date, and so do their badges on shopping lists.
8. **Store editions on demand:** a Kaufland store edition is only read once at least one user has that store in range.
9. **Failure:** if a source changes or blocks access, the founder is alerted and users of that chain see "flyer not updated, scan it yourself".

### Location (all without AI)

- **Shop positions** come from OpenStreetMap, which is open data and covers all EU countries, so no per-country deal is needed. Users can add or correct a shop.
- **User position** comes from the phone, asked for only when the user first opens offers. If the user declines, a typed postcode or town works instead.
- **Privacy:** distance is calculated on the phone; the server receives at most a rounded area (postcode level) to pick regional flyers. No location history is kept.

## 4. Where product barcodes (EANs) come from

The barcode on a pack is an EAN (also called GTIN): 13 digits, or 8 on small packs. Scanning it is the fastest way to stock the fridge, but only if the app can turn the number into a product.

### How one scan is resolved (without AI)

The app tries these in order and stops at the first hit:

1. **This household's own products:** anything scanned here before. Works offline.
2. **Coolinary's shared catalogue:** every product any user has confirmed.
3. **Open Food Facts:** the open food product database.
4. **The official GS1 register,** where access has been agreed: name, brand, and brand owner.
5. **No hit:** the user names the product once, or photographs the label and AI reads it (per request). It joins the shared catalogue and counts as confirmed once a second household accepts it unchanged.

After a hit, the product is matched to a catalogue ingredient, and the storage place and use-by date are suggested from the product type.

### Sources of EAN data

| Source | What it gives | Terms and cost | Role |
|---|---|---|---|
| Users' own scans | The code plus a name typed once | Free | Core; fills local gaps such as regional brands and shop own-brands |
| Open Food Facts | Name, brand, size, category, photo. About 429,000 products are listed for Germany (5 October 2026); the Poland figure could not be read and must be measured | Free under the Open Database Licence: credit is required, and an adapted database used publicly must be offered under the same licence. Live lookups are allowed when one call is one real scan; bulk downloads exist for caching. Product additions should be sent back | First outside source |
| GS1 registers: eProdukty (GS1 Polska), Verified by GS1 | Official name, brand, category, and brand owner for a registered code | The public Verified by GS1 page allows 30 lookups a day; API access goes through a GS1 member organisation (terms and price to confirm) | Second outside source; confirms a code is genuine |
| Paid lookup services, for example EAN-Search | Name and category for codes the others miss | Paid (price to confirm) | Optional fallback |
| 2D codes on the pack (GS1 QR or DataMatrix) | Product code, use-by date, and batch in one scan | Free to read, still rare. GS1's "Sunrise 2027" goal is for shop tills to read them by the end of 2027 | Read from release 1; removes typing the date |
| Receipts | Product names as printed, usually without the EAN | AI reads the photo (per request) | Lines are matched to products this household scanned before at that shop |
| Shop loyalty apps (digital receipts) | The same names as a paper receipt | No public access; the user can share the receipt image or PDF | Same as receipts |
| Flyers | No EANs | Not applicable | Offers are matched to ingredients by name, not by code |

### Codes that need special handling

- **Codes starting with 20 to 29** are shop-internal (weighed cheese, meat, bakery). No register holds them; the digits carry a price or a weight under rules that differ by country and chain. The app reads the weight where it knows the chain's pattern and otherwise asks for a name once per shop.
- **Loose fruit and vegetables** have no code: pick from a list or use a photo.
- **Multipacks** have their own code, different from the single item; the catalogue stores how many units a pack holds.
- **8-digit codes** and 12-digit codes on imported goods are read the same way as 13-digit ones.

### Ways to stock the fridge, fastest first

1. **Tick in the shop:** an item ticked on the shopping list moves into the fridge.
2. **Unpack mode:** scan every barcode in the bag in a row, then confirm the whole list once.
3. **Receipt photo** (AI per request).
4. **Shelf photo** (AI per request).
5. **Typing,** with autocomplete.

### Open points

- Licence check before launch: how Coolinary's shared catalogue may be combined with Open Food Facts data, given the share-alike rule.
- GS1 access terms and price for Poland and Germany.
- The share of products recognised on first scan in each country (verification item 14).

## 5. Where nutrition values come from

### How a dish's calories are calculated (without AI)

1. Each catalogue ingredient is matched once to an entry in a nutrition table. The entry holds values per 100 g: energy, protein, carbohydrate, sugars, fat, saturates, fibre, and salt.
2. Amounts in pieces or spoons ("1 onion", "2 tbsp oil") become grams through the catalogue's typical weights.
3. Ingredients are counted as weighed in the recipe (dry pasta, raw meat). For deep frying, only the oil the food absorbs is counted.
4. When the household cooks with a product it scanned, and the pack's values are known, those values replace the table's. Packs sold in the EU must state energy, fat, saturates, carbohydrate, sugars, protein, and salt per 100 g.
5. The total is divided by the dish's portions, rounded to the nearest 10 kcal, and shown as "about".
6. If any ingredient has no match, the dish shows "calories incomplete" and names the missing ingredient, rather than a number that is too low.

### Sources of nutrition values

| Source | What it gives | Terms and cost | Role |
|---|---|---|---|
| Ciqual (ANSES, France) | 3,484 foods, 74 nutrients each (update of November 2025) | Free, published as open data; exact licence and credit wording to confirm | First source for plain ingredients |
| Bundeslebensmittelschlüssel 4.0 (Max Rubner-Institut, Germany) | About 7,140 foods and dishes, up to 138 nutrients each, cooked foods included | Free since December 2025; terms for a commercial app to confirm | German foods, and a second source where Ciqual has no entry |
| USDA FoodData Central | US plain and branded foods | Public domain (CC0); credit requested, not required | Fills gaps such as Asian sauces and spices |
| Open Food Facts | Pack values of branded products | Open Database Licence (see section 4) | Products the household scanned |
| Polish tables (NIZP PZH-PIB, "Tabele składu i wartości odżywczej żywności") | Polish foods and dishes | Published as a book; reuse terms not found | Only once a licence is agreed. Until then, Polish ingredients use the closest Ciqual or BLS entry |

### Open points

- Licence check of Ciqual and BLS for use in a paid app, and the credit each one requires.
- Whether NIZP PZH-PIB licenses its tables for app use, and at what price.
- A dietitian checks the calorie values of the 200 starter dishes before launch.
- Before release 2, a paediatric dietitian reviews the children's goals: whether "grow into a healthier weight" is safe as described, and the lower limit per age for numbers set by a doctor or dietitian.
- Meal shares by country: Polish homes eat the main hot meal at midday, and German homes often have a cold supper. The default 25/40/25/10 split may need a version per country.

## 6. Free and Plus (mirrored in business model section 3)

| | Free | Plus |
|---|---|---|
| Without AI | Fridge list, barcode, library matches, cuisine and tool filters, allergen filter, calories and nutrients per dish, one personal calorie goal, shopping list with missing items, offers from the 2 nearest shops | Adds: expiry reminders, all shops in the travel range, the trip plan across shops, near-shop reminders, shared household (5 people), per-person diet profiles, calorie goals and portions for every household member (children included), savings report |
| AI per request | 2 day-plans and 2 photo scans per week | Unlimited: scans, day and week plans, family plans to everyone's goals with the fewest cooking sessions, swaps, tool rewrites, plan from deals, own-flyer scan, cheapest top-up |

Changes from the first version of the business model: the allergen filter and the basic missing-items list move to Free (safety, and shopping lists are free in every rival app); flyer and location features are new. Calories per dish are free, as in most rival apps; goals for several people join the per-person profiles in Plus.

## 7. Release order

| Release | Contents | Reason |
|---|---|---|
| **1. Launch** | Kitchen profile (cuisines, tools, diet); fridge by typing, barcode (with unpack mode and naming of unknown products), and photo; AI day plan; starter library of about 200 checked dishes in 5 cuisines (Polish, German, Italian, Asian, Mediterranean) with non-AI matching; calories and nutrients per dish; cooking mode; shopping list with missing items; expiry reminders; paywall | The core loop, and enough to charge for. Calories are expected in any recipe app, and are cheaper to calculate for the starter library before launch than to add later |
| **2. Flyers and family** | Own-flyer scan (all countries); synced flyers for the launch chains (Biedronka, Lidl, Kaufland in Poland; Lidl, Aldi Nord, Aldi Süd, Kaufland in Germany); nearest-shops-first ordering with travel range and postcode fallback; store edition picked by location; offer browsing; offer badges on the list; trip plan across shops within the walking limit; plan from deals; shared household; personal and family calorie goals with one-dish portions; profiles without accounts for children | A weekly reason to return; sets Coolinary apart from scan-only apps. Family goals need the shared household, so they arrive with it |
| **3. Depth** | Week plan with leftovers and the fewest cooking sessions (cook once, eat twice); cheapest top-up; receipt and voice input; restyle a dish; next chains (Dino, Netto, Carrefour, Auchan, Aldi in Poland; Penny in Germany); near-shop reminders; licensed flyer feed (Rewe, Edeka, Netto) and more countries; shareable savings card | Retention and the second revenue line |

Before building release 2, add two variants to the concierge test in the business model's validation phase: "plan from this week's flyer", and "one dish, a portion for each person" for families whose members have different needs.

## 8. Risks

| Risk | Response |
|---|---|
| Shelf prices are not in flyers | Until receipt prices build up, only flyer offers separate the shops; totals say "about" and the plan never claims more than it knows |
| A barcode is not recognised | The user names the product once and it is known to everyone afterwards; the hit rate is measured per country before launch |
| The open product database must be shared alike | Legal check before launch; product additions are sent back; household data is kept separate from the product catalogue |
| A flyer price is read wrongly | Link to the source flyer on every offer; a chain goes live only at ≥ 97% correct prices and dates on a 100-offer sample; weekly spot-check of 20 |
| A retailer objects | Remove the chain within a day; user-supplied flyers keep working; move to the licensed feed |
| Flyva already plans meals from German flyers | Coolinary stays fridge-first: deals only fill the gaps left after using what is at home |
| Prices that need a loyalty app or card | Shown as a condition on the offer |
| Flyers differ per store or region | Edition picked by location; store editions read only where there are users; fully per-store chains (Rewe, Edeka) wait for the licensed feed |
| A chain changes its flyer page or blocks access | Daily check raises an alert; that chain falls back to user-supplied flyers until fixed |
| Users refuse to share location | Postcode or town entry gives the same ordering; nothing else in the app depends on location |
| Location is sensitive personal data | Distance calculated on the phone, only a rounded area sent, no history kept; near-shop reminders are opt-in |
| A shop is missing or misplaced on the map | Users can add or correct their shop |
| AI suggests a dish with a missing tool or an allergen | The non-AI check in rule 2 blocks it before display |
| A calorie value is wrong | Values come from official tables or the pack, shown as "about"; a dish with an unmatched ingredient says "calories incomplete"; a dietitian checks the starter library |
| Goals are read as medical advice | Estimates are "a starting point"; no goals tied to an illness and no health claims, so the app stays a lifestyle app outside the EU medical device rules (MDCG 2019-11); filters state numbers, not "healthy" or "light" |
| Calorie goals harm someone with an eating disorder | No adult target under 1,200 kcal; no weight history; calories can be hidden, and are hidden by default for children |
| A child is put on a diet or feels singled out | No formula-based "lose" goal under 18; only a parent sets a child's goal; a child's target never goes below the EFSA value for their age at the lowest activity level, unless a doctor or dietitian sets it, and even then not below a lower limit per age; the whole family eats lighter together; the child sees portions only |
| An overweight child's family gets no help | "Grow into a healthier weight" and "set by our doctor or dietitian" goals, with "eat lighter as a family" |
| Weight and goals are health data under the GDPR | Separate explicit consent; kept on the phone; only the daily target reaches the server; other adults see portions, not goals |
| Family needs are too far apart for one dish | Portions from ½ to 2, no-cook extras, and "base plus toppings" dishes; a second dish only when an allergy or diet forces it |
| Onboarding feels long | Presets and skippable defaults; people and goals are asked later, not at first start |

## Verification

Checks on each finished release:

1. **AI off:** with AI disabled and 10 items in the fridge, the user gets at least 3 library matches per meal, filtered by cuisine and tools, and can build a shopping list.
2. **Tools:** with the "microwave + kettle" preset, 50 composed dishes need no other tool.
3. **Allergens:** with "no gluten, no nuts", 50 composed dishes contain zero violations.
4. **Cuisine:** for each launch cuisine, a day plan from a fixed test fridge is either in that cuisine or states which items are missing.
5. **Flyer accuracy:** the 97% sample check passes for each chain before it goes live.
6. **Flyer to list:** five offers added from a flyer appear on the list with shop, price, and end date; the badge disappears when the offer ends.
7. **Sync:** for each launch chain over two weeks, every new flyer on the chain's list page shows up in the app before its start date, expired offers are gone on their end date, and a simulated source failure triggers the alert and the fallback message.
8. **Edition:** test positions in Hamburg and Munich get Aldi Nord and Aldi Süd offers respectively; two Kaufland stores with different editions show different offers.
9. **Nearest first:** from a fixed test position, offers appear in order of shop distance and nothing outside the travel range is shown; changing the range changes the set.
10. **No location:** with permission denied, entering a postcode gives the same ordering, and the rest of the app works unchanged.
11. **Cost:** over the first 100 users, AI cost stays within €0.10 a month per Free user and 15% of net revenue per Plus household.
12. **List making:** for a fixed test plan and test fridge, the proposed list holds exactly the missing ingredients, each once, with amounts merged across dishes; nothing already at home appears; ticking an item in the shop adds it to the fridge.
13. **Trip plan:** with the worked example in section 2, the plan is Lidl only (10.65) at a €2 minimum and Lidl then Kaufland (9.15) at a €1 minimum; with walking set to none only single shops are offered; with the walk limit under 900 m the two-shop plan disappears.
14. **Barcode:** 50 packaged products from one shop in Poland and 50 from one in Germany are scanned and the share recognised on first scan is recorded (target at least 80%); an unknown product named once is recognised on the next scan from a second account; a weighed-counter code, an 8-digit code, and a multipack are handled.
15. **Linking:** a second test account joins by QR, by code, and by link; a used or expired invite is refused; both phones then show the same fridge and list; after the owner removes the second account, it no longer sees them.
16. **Calories:** for 20 starter dishes, calories per portion are within 10% of a dietitian's calculation; a dish with an unmatched ingredient shows "calories incomplete", never a number; with AI off, calories and portions still work.
17. **Portions:** with the worked example in "Planning for each person", the pot is 6 portions (Adam 1½, Ewa 1½, Ola 1, Kuba 2) and only Kuba gets an extra; marking Ola away for dinner makes the pot 5 portions and the shopping list follows.
18. **Family week:** for a test family of four with different goals and a fixed fridge, a 7-day plan set to "every other day" has at most 4 cooking sessions and one dish per meal (unless an allergy or diet forces a second), with each person's days within 10% of target and their week within 5%.
19. **Goal limits:** an adult target under 1,200 kcal is refused by every route, including "my own number". For a profile under 18: no formula-based "lose" goal; "grow into a healthier weight" never gives less than the EFSA value for the child's age and sex at the lowest activity level; a doctor's number under the age limit is refused; a teenager's own account cannot change their goal; the child's view shows no calories and no goal; turning on "eat lighter as a family" changes the extras for everyone, not only the child.
20. **Health data:** network traffic shows no age, height, or weight leaving the phone; a second adult in the household sees portions but not the other adult's weight or goal; deleting a goal removes the target from the server.

## Sources

- Flyva, "Cook what's on sale near you": https://www.producthunt.com/p/flyva/flyva
- Offerista / ShopFully network and Publitas cooperation: https://www.publitas.com/blog/the-future-of-the-brochure-offerista-and-publitas-launch-europe-wide-cooperation/
- Offerista Group (now ShopFully): https://app.dealroom.co/companies/offerista_group
- Polish flyer and shopping-list apps (Blix, Moja Gazetka): https://nano.komputronik.pl/n/aplikacja-do-listy-zakupow-ranking-top-5/
- Flyer pages: [Biedronka](https://www.biedronka.pl/pl/gazetki), [Lidl PL](https://www.lidl.pl/c/nasze-gazetki/s10008614), [Kaufland PL](https://sklep.kaufland.pl/gazeta-reklamowa.html), [Lidl DE](https://www.lidl.de/c/online-prospekte/s10005610), [Aldi Süd](https://www.aldi-sued.de/de/angebote/prospekte.html), [Aldi Nord](https://www.aldi-nord.de/prospekte.html), [Kaufland DE](https://filiale.kaufland.de/prospekte.html), [Dino](https://marketdino.pl/gazetki/lista/standard), [Netto PL](https://netto.pl/gazetka-netto/), [Carrefour PL](https://www.carrefour.pl/gazetka-handlowa), [Auchan PL](https://www.auchan.pl/pl/gazetki), [Aldi PL](https://www.aldi.pl/gazetki.html), [Penny](https://www.penny.de/angebote)
- Open Food Facts, data reuse and licence: https://world.openfoodfacts.org/data
- GS1 Polska eProdukty: https://itwiz.pl/eprodukty-globalna-baza-kodow-ean-i-baza-produktowa-na-wyciagniecie-reki/
- Verified by GS1 (30 free lookups a day, API through member organisations): https://orders.gs1.de/en/verified-by-gs1-will-replace-gepir
- GS1 Sunrise 2027 (2D codes at the till): https://www.dynamsoft.com/blog/insights/gs1-sunrise-2027-transition-to-2d-barcodes-for-global-commerce/
- Shop-internal codes with prefixes 20 to 29: https://gs1.se/en/guides/how-to-guides/barcode-label-items-of-varying-weight/
- GS1 prefix list (952 is reserved for examples): https://www.gs1.org/prefixes
- Ryanair v PR Aviation, contractual limits on data reuse: https://www.pinsentmasons.com/out-law/news/cjeu-case-on-screen-scraping-has-potential-to-affect-business-models-says-expert
- Ciqual 2025 update (3,484 foods, 74 nutrients): https://www.anses.fr/en/content/new-enhanced-and-more-representative-version-ciqual-table
- Bundeslebensmittelschlüssel 4.0 free of charge (7,140 foods, up to 138 nutrients): https://heise.de/-11123877
- USDA FoodData Central, public domain (CC0): https://fdc.nal.usda.gov/api-guide
- Polish food composition tables (NIZP PZH-PIB): https://www.pzh.gov.pl/wp-content/uploads/2021/07/OKLADKA-WERSJA-PELNA.pdf
- Mandatory nutrition declaration on EU packs: https://food.ec.europa.eu/food-safety/labelling-and-nutrition/food-information-consumers-legislation/nutrition-labelling_en
- EFSA dietary reference values for energy (adults and children by age, sex, and activity): https://efsa.europa.eu/en/efsajournal/pub/3005
- WHO Europe COSI, round 6 (2022–2024): 25% of children aged 7–9 with overweight, 11% with obesity: https://www.who.int/europe/news/item/29-11-2024-new-who-europe-fact-sheet-highlights-worrying-post-covid-trends-in-childhood-obesity
- NICE PH47, family-based weight management for children; all family members encouraged to eat healthily whatever their weight (now under NG246): https://nice.org.uk/guidance/ph47/resources/weight-management-lifestyle-services-for-overweight-or-obese-children-and-young-people-1996362978757
- AAP 2023 guideline: family-based, intensive lifestyle treatment for children with obesity: https://www.eurekalert.org/news-releases/976540
- Kurbo by WW, backlash over a weight-loss app for children (2019): https://edr.iaedpfoundation.com/update-new-weight-watchers-anti-obesity-app-for-children-draws-criticisms/
- NHS calorie guidance (a cut of about 600 kcal a day is safe; Coolinary uses a gentler 500): https://www.nhs.uk/better-health/lose-weight/calorie-counting/
- MDCG 2019-11, lifestyle and well-being software is not a medical device: https://health.ec.europa.eu/document/download/b865d8e9-081a-4601-a91a-f120321c0491_en
- Eat This Much, family meals (one plan scaled up; members' different targets are not met at the same time): https://blog.eatthismuch.com/?p=413
