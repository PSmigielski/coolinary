# Coolinary — Business Model Plan

Scope: business model only. No product, tech, or design decisions.
Decisions already made: **EU multi-country launch**, **solo / bootstrapped**.

## Context

EU households throw away about 72 kg of food per person per year, 54% of all EU food waste (Eurostat, 2022). Since October 2025 member states are legally bound to cut retail + consumption food waste by 30% per capita by 2030. At the same time, "what do we eat today?" is a daily decision most households dislike making.

Coolinary turns the contents of the fridge into a full day of meals: **breakfast, dinner, and supper**, built from what is already there, oldest items first.

Two facts shape the whole model:

1. **"Photo of fridge → recipe" is already a commodity.** A dozen indie apps sell it for $18–70/yr, Samsung Food+ sells it for $59.99/yr, and general AI assistants do it for free. Coolinary cannot charge for the scan. It has to charge for the outcome: the day is planned, money is saved, nothing is thrown out.
2. **Solo + multi-country only works because AI makes localisation cheap.** The scarce resources are the founder's attention and marketing reach, so the model must run on organic growth and roll out by language in waves.

## 1. Positioning

**"One look in your fridge, three meals sorted."**

| Rivals sell | Coolinary sells |
|---|---|
| Single recipes from a scan | A whole day composed from the same stock, so breakfast does not use up what dinner needs |
| English/US-first content | EU-native: local language, metric units, local staples and meal rhythm (PL obiad at midday, DE cold Abendbrot, late southern suppers) |
| One-off answers (chat assistants) | A persistent, shared household fridge that remembers stock, tastes, and expiry dates |
| "Cool AI" | Visible savings: "€38 and 4 kg saved this month" |
| Calorie planners built for one person | One shared dish per meal, a portion sized for each family member, and as few cooking sessions as possible |

## 2. Customers

| Priority | Segment | Role in the model |
|---|---|---|
| 1 | Households of 2–4 that cook at home 4+ days a week and feel grocery prices | Core payers |
| 2 | Students, singles, flat-sharers | Free tier, word of mouth, content reach |
| 3 | Diet-specific cooks (vegetarian, gluten-free, high-protein) | Higher willingness to pay; later upsell |
| 4 | Municipalities, waste operators, retailers, food brands | B2B revenue from Phase 3 |

## 3. Offer and pricing

One paid plan, sold per household, annual-first (cash up front, lower churn).

| | Free | Plus (whole household) |
|---|---|---|
| Included | 2 AI day-plans and 2 photo scans per week, 1 person. Everything that runs without AI: fridge list, barcode, recipe-library matches, cuisine and tool filters, allergen filter, calories per dish, one personal calorie goal, shopping list with missing items, offers from the 2 nearest shops | Unlimited AI (scans, day and week plans, family plans to everyone's goals with the fewest cooking sessions, swaps, plan from deals, own-flyer scan, cheapest top-up); expiry reminders; offers from all shops in range with cheapest-shop comparison; shared fridge for up to 5 people; per-person diet profiles and calorie goals, children included; monthly savings report |
| Tier A: DE, FR, NL, AT, BE, IE, Nordics | €0 | €4.99/mo or €34.99/yr |
| Tier B: PL, CZ, SK, HU, RO (local currency) | €0 | ≈ €2.99/mo or ≈ €19.99/yr |

- The full feature split (with and without AI) is in `docs/features.md`.
- Southern Europe (ES, IT, PT): test both tiers.
- Paywall: first scan and first full day-plan are free with no account, then a paywall with a 7-day trial on the annual plan; dismissing it drops the user to the Free tier. Test this against a hard paywall early: the 2026 RevenueCat benchmark is 10.7% download-to-paid for hard paywalls vs 2.1% for freemium.
- Price anchors: Fridgi £34.99/yr, Samsung Food+ $59.99/yr, SnapChef $59.99/yr, Primo $17.99/yr, Spoilr $19.99/yr.

## 4. Revenue streams

| Stream | Starts | Expected share | Notes |
|---|---|---|---|
| Plus subscription | Launch | 85–100% in year 1 | The only stream the plan depends on |
| Grocery affiliate on "missing items" lists | Phase 3 | 5–10% | Via pan-EU affiliate networks, so no per-country deals |
| Sponsored "rescue recipes" from food brands, clearly labelled | Phase 3, at ≥ 50k monthly users | 5–10% | Sold per campaign |
| Offers network (for example ShopFully) paying for flyer views | Phase 3 | Upside | Depends on traffic; also the licensed flyer source for per-store chains |
| B2B / B2G licences | Phase 3 | Upside | Co-branded resident or customer campaigns for cities, waste operators, and retailers that must hit the 2030 target; bulk Plus licences or a flat campaign fee |

Ruled out: banner ad networks, selling personal data, hardware.

## 5. Unit economics (base case, all to be validated)

| Item | Value |
|---|---|
| Blended annual price | €28 |
| After VAT (≈ 21%) | €23.1 |
| After 15% store fee (small-business programmes) | €19.7 |
| AI and hosting cost per payer per year (budget cap) | ≤ €3.5 |
| **Contribution per payer per year** | **≈ €16** |
| Expected payer lifetime (assumption) | ≈ 1.5 years, so LTV ≈ €24 |
| Install → paid (assumption) | 2.5%, so ≈ €0.60 LTV per install |

- **Target: €3,000 net per month** ("ramen profitable") needs ≈ 2,250 active payers, which is ≈ 90,000 installs at 2.5%.
- **Paid advertising does not work at these numbers.** A healthy 3:1 ratio allows about €0.20 per install; EU installs cost several times that. Growth has to be organic until install → paid reaches ≥ 5%.

## 6. Cost structure

- **Variable:** AI usage per scan and per plan (main cost), store fees, payment/VAT handling.
- **Fixed, kept under ≈ €200/month before revenue:** developer accounts, hosting, subscription and analytics tools, accounting.
- **Fixed from the flyer release:** weekly AI reading of each covered chain's flyers. This cost depends on the number of chains and editions, not on the number of users.
- **One-off per language:** native-speaker check of the app text and store listing (a few hundred euros).
- **Guardrails:** a free user must cost ≤ €0.10/month (enforced by the 2-plans-per-week cap); AI cost per payer ≤ 15% of net revenue.

## 7. Go-to-market: multi-country, one person

**Rollout in language waves**
- Wave 1: English (available EU-wide) + German + Polish. Polish is assumed to be the home market for hands-on interviews; swap if wrong.
- Wave 2: French, Spanish, Italian.
- Wave 3: Dutch, Nordics, Czech, others on demand.
- Rule: after each wave, put most effort into the two countries with the best install → paid rate. A language that stays under 1% after 3 months gets no further marketing.

**Channels, in order**
1. Store listings localised per language. Non-English stores are under-served by US-first rivals.
2. Short video: the "one fridge → three meals" format, one clip a day, plus creators paid by revenue share instead of fees.
3. Built-in sharing: household invites (shared fridge) and a shareable monthly savings card.
4. Communities and press: zero-waste and budget-cooking groups, student unions, the food-waste angle around 29 September (UN food-waste awareness day) and the 2030 targets.
5. Owned audience: waitlist and email list per language, to reduce dependence on the app stores.

## 8. Competition

| Type | Examples | How Coolinary competes |
|---|---|---|
| Indie fridge-scan apps | Fridgi, SnapChef, Scrappy AI, Primo, Spoilr | Full-day planning, EU languages, household plan at a mid-range price |
| Platform apps | Samsung Food+ | Works with any fridge, half the price, local cuisine |
| Grocer-funded planners | Jow (free; earns from grocers and retail media) | Uses what you already have instead of selling a new basket |
| Deals-first meal planners | Flyva (Germany) | Fridge-first: offers only fill the gaps left after using what is at home |
| Flyer and shopping-list apps | Blix, Moja Gazetka (Poland) | They list offers; Coolinary turns offers plus the fridge into meals |
| Calorie-based meal planners | Eat This Much | Its family setting scales one person's plan and does not aim at each member's target at once; Coolinary sizes each person's portion from one shared pot |
| General AI assistants | ChatGPT, Gemini, Claude | Remembers the household's stock, expiry dates, and tastes; shows savings |

Defensibility comes from the stored household fridge and taste history, the daily habit, and store rankings in many languages. The scan itself gives none.

## 9. Roadmap with gates

| Phase | Goal | Gate to move on |
|---|---|---|
| 0. Validate (≈ 4 weeks, no product) | Prove the pain and willingness to pay | ≥ 300 waitlist sign-ups across ≥ 2 countries; ≥ 30% of concierge testers come back 3+ times in a week; ≥ 10 pre-paid founding annual plans |
| 1. Paid launch, Wave 1 | First revenue | 100 payers; install → paid ≥ 2%; AI cost within the 15% cap; store rating ≥ 4.5 |
| 2. Growth loop, Wave 2 | One repeatable organic channel | 1,000 payers; ≥ 80% of installs organic; 30-day retention ≥ 12% |
| 3. Diversify | Ramen profitability and a second revenue stream | ≈ 2,250 payers (€3k net/month); affiliate or B2B ≥ 10% of revenue |

Stop or pivot rules: if Phase 0 misses its gates, narrow the audience (for example families with children, or one diet) and retest before building. If six months after launch install → paid is under 1% and 30-day retention under 8%, move to B2B licensing or stop.

## 10. Risks

| Risk | Response |
|---|---|
| The core feature becomes free everywhere | Sell the outcome and the household memory; keep the savings report central |
| People will not pay for recipes; high churn | Annual-first pricing; savings shown to exceed the price; household plan |
| Free users drive up AI costs | Hard weekly cap on Free; track cost per active user weekly |
| One founder spread across too many countries | Language waves and the "top two countries" rule |
| Allergen or food-safety complaint | Allergen profiles, conservative expiry guidance, clear AI-generated label, no health or medical claims |
| Dependence on app stores | Web landing pages and email list per language |
| Privacy (fridge photos, household data, location, body weight and goals) | Collect the minimum, EU hosting, distance to shops calculated on the phone with no location history, weight and goals kept on the phone with only the daily calorie target synced, sell privacy as a feature |
| Calorie goals are read as medical advice or harm someone | Lifestyle wording only, no adult target under 1,200 kcal, no weight history; a child's goal is set by a parent and is never a formula-based calorie cut: the whole family eats lighter together and the child is not singled out (details in `docs/features.md`) |

## Verification (Phase 0, before anything is built)

1. **Interviews:** 15–20 home cooks across Poland, Germany, and one more country. Confirm the daily-decision pain and what they pay for today.
2. **Waitlist page in 3 languages** with two price points shown at random; measure sign-up and "reserve at this price" clicks.
3. **Concierge test:** 30–50 people send a fridge photo by chat and get a breakfast/dinner/supper plan back the same day. Measure repeat use within a week and whether meals were actually cooked. Give half of them a "plan from this week's flyer" variant (their fridge plus current offers at their nearest shop) and compare repeat use. Give families whose members have different needs (for example a dieting adult and a growing teenager) one dish per meal with a portion for each person, and record whether they cooked once per meal.
4. **Pre-sale:** offer a founding annual plan at €19 to testers and the waitlist.
5. **Model check:** put real numbers from steps 2–4 into the section 5 table and re-check the €3k/month target and the Phase 0 gates.
6. **Competitor check:** re-verify current prices and languages of the apps in section 8 in the German, Polish, and French stores.

## Sources

- Eurostat, "Food waste: 132 kg per inhabitant in the EU in 2022": https://ec.europa.eu/eurostat/web/products-eurostat-news/w/ddn-20240927-2
- European Commission, food waste reduction targets: https://food.ec.europa.eu/food-safety/food-waste/eu-food-waste-relevant-legislation/food-waste-reduction-targets_en
- RevenueCat, State of Subscription Apps 2026: https://www.revenuecat.com/state-of-subscription-apps
- Fridgi (App Store): https://apps.apple.com/app/id6766045022
- Spoilr (App Store): https://apps.apple.com/us/app/id6760987461
- Scrappy AI (App Store): https://apps.apple.com/app/id6781905536
- Primo: https://appgoblin.info/apps/6747734717
- SnapChef: https://appgoblin.info/apps/6783267721
- Samsung Food+ pricing: https://stuff.co.za/2024/09/06/samsung-food-upgraded-ai-app-pay-for/
- Jow (Northzone investment note): https://northzone.com/insights/our-investment-in-jow-revolutionising-the-grocery-shopping-experience-in-europe-and-the-us
- Eat This Much, family and couple meal planning: https://blog.eatthismuch.com/?p=413
