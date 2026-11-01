# Build requirements for a credible calorie-tracking food diary inside Coolinary (Poland and Germany, as of October 2026)

Notes conventions:
- **Confirmed** = has a URL. **Estimate** = my arithmetic, with assumptions written out. **Snippet** = taken from a search-result summary and not opened in full.
- USD→EUR: I use an **assumed rate of €0.86 per $1**. It was not verified for October 2026, so re-check before quoting. Figures are rounded.
- Research date: 2026-10-30. Anything older than 2024 is flagged.

---

## 1. Feature baseline users expect from a calorie tracker in 2026: table stakes vs differentiators

### Takeaway
Users expect a diary split by meals, fast search over a large local database, barcode scanning, quick-add, history and a weight log. In 2026 AI photo or text logging is now a paid-tier feature at the market leaders (Fitatu Premium+AI, MyFitnessPal Meal Scan plus the Cal AI purchase), so it has moved from novelty to expected premium feature. Coolinary's real differentiator would be "the plan already knows what you ate": one-tap logging of a planned, cooked meal with portions taken from the recipe.

### Cited Findings
- **Fitatu (PL market leader).**
  - Says its database holds ">1 million products and dishes": basic foods, ready meals, store brands, restaurant and chain dishes — [Fitatu help](https://www.fitatu.com/pl/help/co-to-jest-baza-produktow-i-potraw-fitatu).
  - Says the Polish database alone holds ">500,000" items — [Fitatu help](https://www.fitatu.com/pl/help/czy-fitatu-posiada-najwieksza-baze-polskich-produktow-i-potraw).
  - Premium users can filter search by kcal and macro ranges and pick which language database to search (PL, EN, NL, FR, IT, DE, ES) — same Fitatu help pages (snippet).
- **Fitatu AI.**
  - "Szacowanie kalorii AI ze zdjęcia" (calorie estimate from a photo) is available only in **Premium+AI**, not in plain Premium. The photo result can be refined with a text description — [Fitatu help](https://www.fitatu.com/pl/help/czym-jest-funkcja-%22szacowanie-kalorii-ai-ze-zdjecia%22).
  - Premium+AI also includes AI text entry — [Fitatu help: AI features](https://www.fitatu.com/pl/help/jakie-funkcje-ai-oferuje-aplikacja-fitatu).
  - Versions so far: "Fitatu AI 2.0 Acai" and a later "AI 3.0 Goji" — [Fitatu blog](https://www.fitatu.com/blog/?p=14101).
  - Two Fitatu pages disagree on whether gallery photos are supported (snippet).
- **Fitatu Premium+AI price.** Sources conflict:
  - 49.99 zł/month, or 169.99 zł/year on promotion — [Telepolis](https://www.telepolis.pl/tech/aplikacje/fitatu-ai-kalorie-ze-zdjecia-posilku).
  - 64.99 zł/month or 189.99 zł/year — [Tabletowo](https://www.tabletowo.pl/?p=960186).
  - Neither article is clearly dated (Fitatu's AI launch coverage is from around March 2025: [Spidersweb](https://www.spidersweb.pl/2025/03/fitatu-aplikacja-liczenie-kalorii-sztuczna-inteligencja-smartfon.html)).
- **MyFitnessPal.**
  - Database: "20 million foods, 68,500 brands and meals from 380+ restaurant chains".
  - Has a photo "Meal Scan" feature; users can correct detected items (e.g. three pickles, not two).
  - Bought Cal AI (deal closed December 2025, announced March 2026). The CEO describes Cal AI as built "for speed rather than precision".
  - All from [TechCrunch, 2 Mar 2026](https://techcrunch.com/2026/03/02/myfitnesspal-has-acquired-cal-ai-the-viral-calorie-app-built-by-teens).
- **Yazio.** The figures conflict:
  - The App Store listing claims a calorie table of more than 2 million foods — [App Store PL](https://apps.apple.com/pl/app/id946099227?l=en-gb) (snippet).
  - A review site says about 2.5M — [calorie-trackers.com](https://calorie-trackers.com/reviews/yazio/).
  - A competitor (Nutrola) says 4M+ — [Nutrola](https://nutrola.app/en/blog/calorie-tracker-database-size-comparison-2026).
  - Yazio is described as strongest for DACH products — same sources (marketing-adjacent).
  - AI photo logging: competitor blogs disagree on whether Yazio has it. VitaCal (April 2026) says it is behind PRO — [VitaCal](https://vitacal.pages.dev/compare/vitacal-vs-yazio).
  - Nutrola says Yazio has no voice logging in 2026 — [Nutrola](https://nutrola.app/en/blog/why-does-yazio-not-have-voice-logging). This is a competitor source, not neutral.
  - PRO price: about $6.99/month or $29.99/year (competitor source) — [Nutrola](https://nutrola.app/en/blog/is-yazio-still-good-2026).
- **MacroFactor** (bootstrapped, launched 2021) is the reference for a "smart" tracker: adaptive targets, coaching. It won Google Play Best of 2024 — [MacroFactor](https://macrofactor.com/google-award/); [MacroFactor team](https://macrofactor.com/team/).

### Inferences
Tiering below is my synthesis from the leaders' feature sets as commonly listed. I did not re-verify each feature in every app this session.

| Feature | Tier | Coolinary today |
|---|---|---|
| Diary by meal slot (breakfast…snacks), daily kcal and macro totals vs target | **Table stakes** | Partly: per-dish kcal and macros plus per-person targets exist; no diary |
| Food search over a local-language DB (PL/DE generic + branded) | **Table stakes** | Ingredient catalogue matched to Ciqual/BLS/USDA; branded search limited to own/crowd/OFF/GS1 |
| Barcode scan → log | **Table stakes** | Scanner + lookups exist; needs "log a serving" UI |
| Quick-add kcal (and optional macros) | **Table stakes** (trivial) | New |
| Custom foods and custom recipes | **Table stakes** | Recipes and dishes exist; user custom foods partly (crowd catalogue) |
| Copy meal or previous day; recents and favourites | **Table stakes** (key retention driver) | New |
| Macro goals | **Table stakes** | Exists via targets (kcal; macro split may need UI) |
| Weight log and chart | **Table stakes** for weight-loss positioning | **Conflicts** with "no weight history" |
| Water tracking | Expected but low value | New (trivial) |
| Activity calories from wearables (Apple Health, Health Connect) | Expected in premium tiers | New |
| Streaks and reminders | Common; engagement lever | New |
| AI photo / text / voice logging | **Premium differentiator, fast becoming expected** (Fitatu Premium+AI, MFP Meal Scan, Cal AI) | New; **conflicts** with "AI never states calorie numbers" |
| Micronutrients, fasting, adaptive TDEE (MacroFactor-style) | Differentiators | Ciqual/BLS hold micronutrients, so a cheap partial win |
| Offline logging | Expected (silent table stakes; users notice when it is missing) | Depends on current sync architecture |

- The Coolinary-specific "killer feature" is logging from the plan: "mark planned meal eaten" plus portion multiplier. No competitor has a fridge-aware plan, so this is both the cheapest and the most defensible feature.

### Gaps
- No official, current feature matrices for Lifesum or Cronometer were checked this session.
- Yazio's official AI-feature status is unconfirmed; only competitor blogs were found.
- No primary data on which features drive retention (e.g. how much "copy previous day" improves D30). Keep that claim as opinion.

---

## 2. The food database problem: sizes, Open Food Facts PL/DE coverage, ODbL, commercial APIs and costs

### Takeaway
The leaders' databases are 10–500× larger than what Open Food Facts offers for Poland:
- Open Food Facts (OFF) in Poland: about 38.5k products, about 26.4k with complete nutrition facts.
- Fitatu's Polish database: >500k items.
- OFF in Germany is much better: about 433k products, about 253k with complete nutrition facts.

Commercial APIs are either US-centric or quote-only for PL/DE (FatSecret Premier). For PL/DE branded coverage, the realistic bootstrapped path is: OFF + generic tables (BLS/Ciqual) + Coolinary's crowd catalogue, with fast "scan → photo of label → submit" fallbacks.

ODbL share-alike matters if OFF data is merged into Coolinary's crowd catalogue and that catalogue is used publicly.

### Cited Findings
- **OFF counts, measured live on 2026-10-30 via the OFF v2 search API** (`/api/v2/search?countries_tags=…&page_size=1`). OFF was intermittently returning 503s; counts retrieved after retries:
  - World total: **4,885,707** products. With `states_tags=en:nutrition-facts-completed`: **3,654,847** (74.8%) — [OFF API](https://world.openfoodfacts.org/api/v2/search?page_size=1&fields=code).
  - **Germany**: **432,640** products; **253,191** nutrition-complete (58.5%) — [OFF API DE](https://world.openfoodfacts.org/api/v2/search?countries_tags=en:germany&page_size=1&fields=code); [OFF API DE complete](https://world.openfoodfacts.org/api/v2/search?countries_tags=en:germany&states_tags=en:nutrition-facts-completed&page_size=1&fields=code).
  - **Poland**: **38,498** products; **26,381** nutrition-complete (68.5%) — [OFF API PL](https://world.openfoodfacts.org/api/v2/search?countries_tags=en:poland&page_size=1&fields=code); [OFF API PL complete](https://world.openfoodfacts.org/api/v2/search?countries_tags=en:poland&states_tags=en:nutrition-facts-completed&page_size=1&fields=code).
  - France (for scale): 1,277,796 — [OFF API FR](https://world.openfoodfacts.org/api/v2/search?countries_tags=en:france&page_size=1&fields=code).
  - Caveats: "countries_tags" means "sold in"; "nutrition-facts-completed" is a contributor-set state, not a verified-accuracy flag.
  - Earlier search results showed much lower per-country figures (e.g. Poland 564 / 176). Those came from Open *Products* Facts and other sister-project pages, not OFF, and should be ignored — [search result page](https://world.openproductsfacts.org/countries).
- **OFF licensing.**
  - The OFF database is under ODbL. Commercial use is allowed. Re-users must credit Open Food Facts with a link, and "derivative works must be shared under the same conditions" — [OFF terms of use](https://world.openfoodfacts.org/terms-of-use).
  - ODbL 1.0 text:
    - A **Derivative Database** includes any adaptation or extraction of a substantial part.
    - Any Derivative Database that is **Publicly Used** must be under ODbL (§4.4), and recipients must be offered a machine-readable copy, free if distributed over the internet (§4.6).
    - **Produced Works** (e.g. a screen showing the kcal of a product) only need a notice (§4.3).
    - Internal use is not public use (§4.5), and keeping the database unmodified inside a **Collective Database** does not trigger share-alike (§4.5).
    - All from [ODbL 1.0](https://opendatacommons.org/licenses/odbl/1-0/).
- **FatSecret Platform API** — [FatSecret API editions](https://platform.fatsecret.com/api-editions):
  - Basic: free, 5,000 calls/day, US only.
  - **Premier Free**: free; for start-ups with <US$1M revenue *and* <US$1M raised; unlimited calls; **US dataset only**; "non-US access available at a discount for eligible users"; attribution required.
  - **Premier**: price "upon request, based on market access", >62 countries, 26 languages, no attribution.
  - Image recognition and NLP are add-ons billed by monthly usage in tiers of 25,000 inputs, with 50% off for start-ups.
  - Default queries hit the US dataset. Non-US localisation requires contacting FatSecret — [FatSecret API](https://platform.fatsecret.com/api) (snippet).
  - FatSecret does run consumer Polish pages ([fatsecret.pl](https://www.fatsecret.pl/kalorie-warto%C5%9Bci%C2%A0od%C5%BCywcze/fit-food)), which suggests a Polish dataset exists. Its API availability and price are unconfirmed.
- **Edamam Food Database API** — [Edamam](https://developer.edamam.com/food-database-api):
  - Plans:
    - Basic: $14/month, 100k calls/month.
    - Core: $69/month, 750k calls/month; Vision pay-as-you-go at $0.015/call.
    - Plus: $299/month, 5M calls/month, 10k Vision requests/month.
    - Unlimited: custom.
  - Attribution is required on all plans.
  - About 0.9–1M foods and 700–790k UPCs. **No mention of EU, Germany or Poland coverage.**
- **Spoonacular** (recipe-centric, US): Cook $29, Culinarian $79, Chef $149/month; Enterprise $300+; 1,500 / 4,500 / 10,000 points/day — [apis.io plans](https://apis.io/plans/spoonacular/spoonacular-plans-pricing/); [apis.io rate limits](https://apis.io/rate-limits/spoonacular/spoonacular-rate-limits/). Third-party aggregator, not the official page; the free tier figure conflicts (50 vs 150 points/day).
- **Nutritionix.**
  - No public list price found; a competitor (Dietly) says pricing is by quote — [dev.to/dietly](https://dev.to/dietly).
  - The API guide is hosted on Syndigo's developer docs — [Syndigo docs](https://docx.syndigo.com/developers/docs/nutritionix-api-guide).
- **Passio Nutrition AI** (SDK and API, token-based) — [Passio cost breakdown](https://www.passio.ai/cost-breakdown):
  - Price per 1M tokens: $2.50 (1–20M), $2.45 (20–100M), $2.40 (100–500M), $2.35 (500M+).
  - An average active user consumes about 250k tokens/month: light 80–100k, power 300–500k+.
  - Image recognition uses 3–5k tokens per image, but Passio's FAQ says 20–30k (conflict on Passio's own site).
  - Plan tiers per a search snippet of passio.ai/features: Starter $99, Growth $599, Pro $2,999/month — [Passio features](https://www.passio.ai/features) (snippet; aggregator sites quote different, older prices).
- **Polish sources.** A search for a Polish branded-product nutrition API with commercial licensing found **nothing usable**. Results were FatSecret consumer pages, company-registry data (Infobel) and lab services — [search results incl. Infobel](https://www.infobelpro.com/companies/poland/food-and-beverage).

### Inferences
- **Coverage arithmetic (estimate).**
  - OFF PL nutrition-complete 26,381 ÷ Fitatu PL claim 500,000 ≈ **5%** of Fitatu's Polish item count.
  - OFF PL total 38,498 ÷ 500,000 ≈ **7.7%**.
  - Fitatu's count includes dishes, generic foods and user entries with duplicates, so true *unique branded* coverage gaps are smaller, but still large.
  - Germany: OFF's 253k nutrition-complete products is a credible launch base (Yazio claims ~2M+ worldwide, not DE-only).
  - **Implication**: a German launch is credible on OFF + BLS. A Polish launch needs a crowd-submission loop from day one.
- **ODbL implication for Coolinary's crowd catalogue.** If OFF records are merged into the crowd catalogue (one table, OFF rows edited or enriched) and served to other users, the merged catalogue is very likely a publicly used Derivative Database. Coolinary would then have to offer it under ODbL and provide a free machine-readable copy (§4.4, §4.6).

  Options:
  - (a) Keep OFF in a separate, unmodified store (Collective Database). Coolinary's own crowd data stays proprietary, and the two are joined only at query time.
  - (b) Contribute crowd submissions back to OFF (makes ODbL moot and improves PL coverage for everyone).
  - (c) Accept ODbL on the crowd catalogue.

  Per-user diaries showing kcal of a product are Produced Works and only need the attribution notice. Interpreting "substantial extraction" is legal judgment, not settled fact.
- **Commercial API cost at 10k and 100k MAU (estimates).**

  Assumptions:
  - 30% of MAU actively diary-log.
  - Each active logger triggers ~200 search/lookup calls/month (≈10/day × 20 logging days, debounced search).
  - The "all MAU log" case is shown as an upper bound.

  | Provider | 10k MAU | 100k MAU | PL/DE fit |
  |---|---|---|---|
  | **OFF self-hosted** (bulk export + Postgres + Typesense/Meilisearch) | Licence €0; infra ≈ €20–100/month (VPS/managed DB guess, unsourced) | ≈ €100–400/month (guess) | DE good, PL weak |
  | **Edamam** | 3,000 × 200 = 600k calls → Core **$69 ≈ €59/month**; if all 10k log, 2.0M → Plus **$299 ≈ €257/month** | 6M–20M calls → above Plus → **Unlimited, custom quote** | No stated EU coverage; attribution required |
  | **FatSecret** | Premier Free **$0** if eligible, but US data only. PL/DE localisation: discounted quote, unknown | Premier quote, unknown | Probably best commercial PL/DE coverage, but price is opaque |
  | **Spoonacular** | $29–149/month | Enterprise $300+ | US/recipe-centric; poor fit |
  | **Passio** (AI + DB) | If 30% are "average" users: 3,000 × 250k = 750M tokens × $2.35 ≈ **$1,763 ≈ €1,516/month**; if all MAU ≈ $5,875 ≈ €5,050 | 30% ≈ $17,625 ≈ €15,160/month; all MAU ≈ $58,750 ≈ €50,500 | Unknown EU coverage |
  | **Nutritionix / Syndigo / 1WorldSync** | Quote only | Quote only | Enterprise PIM-syndication oriented |

- Bootstrapped recommendation: launch on **OFF + BLS/Ciqual + crowd catalogue**. Request a FatSecret Premier Free + PL/DE localisation quote in parallel, as an upgrade path if PL search quality turns out to be the churn driver.

### Gaps
- No FatSecret PL/DE localisation price or per-country product counts. Quote only; the regions page was not found.
- No reliable figure for the branded-product share of Fitatu or Yazio databases.
- No Polish commercial branded-product nutrition DB was found. Not checked: Kcalmar (no API information found), ILEWAZY, Syndigo/1WorldSync pricing, and GS1 Poland/Germany register nutrition-field completeness. GS1 is already in Coolinary's plan; whether those registers hold nutrition panels at scale is unverified.
- OFF data quality: completeness is measured above, accuracy is not. No 2024–2026 study of OFF nutrition-value error rates for PL/DE was found.
- NIZP PZH-PIB table licensing cost not researched (outside the session budget).

---

## 3. AI photo and voice logging: state of the art, accuracy, cost per logged meal

### Takeaway
2025 peer-reviewed studies put general vision-LLM calorie estimates from a photo alone at about **30–36% mean absolute percentage error (MAPE)**, with systematic **under-estimation of large portions**. Error roughly halves when the model is given an ingredient list with amounts (Case 3 in the Nutrients study below: ~53 kcal MAE, ~14% MAPE).

Coolinary already has the recipe and portions for planned meals, so "mark planned meal eaten" beats AI photo logging on accuracy at near-zero cost. AI photo logging is cheap per image with small models: about €0.0003–0.003 per photo on Gemini Flash-Lite / Flash or Claude Haiku 5.5, about €0.0125 on Sonnet-class models.

### Cited Findings
- **Gothenburg study** — [PMC12513282](https://pmc.ncbi.nlm.nih.gov/articles/PMC12513282/):
  - Published in *Current Developments in Nutrition*, 9 Sep 2025.
  - 52 standardised photos; reference = weighed food analysed in Dietist NET (USDA).
  - Energy MAPE: ChatGPT-4o **35.8%**, Claude 3.5 Sonnet **35.8%**, Gemini 1.5 Pro **64.2%**.
  - Mean bias: −10, −18.9 and +65 kcal respectively.
  - Under-estimation grew with portion size (slopes −0.34 to −0.45).
  - Conclusion: "not yet suitable for precise dietary assessment in clinical or athletic populations", but possibly useful for monitoring or screening.
- **ChatGPT-5 study** — [PMC12655113](https://pmc.ncbi.nlm.nih.gov/articles/PMC12655113/):
  - Published in *Nutrients*, 19 Nov 2025; 195 dishes.
  - Energy MAE / MAPE by context:

    | Context | MAE | MAPE |
    |---|---|---|
    | Image only | **123 kcal** | 30.5% |
    | Image + description | **92 kcal** | 24.4% |
    | Image + ingredient list with amounts | **53 kcal** | 13.9% |
    | Ingredient list without image | **66.5 kcal** | 18.1% |

  - Authors say confirmatory evaluation against dietitians is needed before clinical use.
- **ACM BCB 2025** (7 multimodal models) — [ODU digital commons](https://digitalcommons.odu.edu/computerscience_fac_pubs/417) (snippet): Gemini 2.5 Flash carbohydrate MAPE fell from 56.6% to 39.5% with predicted weight. With ground-truth weight it fell to 20.2% (GPT-4.1: 26.8%).
- **Cureus poster** (GPT-5.2, Gemini 3 Flash, Claude Sonnet 4.6) — [Cureus poster](https://www.cureus.com/posters/3217-validity-and-reliability-of-photo-based-ai-estimates-of-meal-calories-and-macronutrients) (snippet; not peer-reviewed): useful for "rough nutritional estimation", "not yet reliable substitutes for validated databases".
- **Market.**
  - Cal AI: >15M downloads, >$30M annual revenue in under two years, 7 employees plus contractors at acquisition — [TechCrunch](https://techcrunch.com/2026/03/02/myfitnesspal-has-acquired-cal-ai-the-viral-calorie-app-built-by-teens).
  - Cal AI was pulled from the App Store in April 2026 for bypassing in-app purchase (Stripe), deceptive paywall pricing display and manipulative repeat offers. It was reinstated after fixes — [Cult of Mac, 22 Apr 2026](https://www.cultofmac.com/news/why-app-store-pulled-cal-ai-app).
  - Fitatu sells photo AI only in Premium+AI — [Fitatu help](https://www.fitatu.com/pl/help/czym-jest-funkcja-%22szacowanie-kalorii-ai-ze-zdjecia%22).
- **Claude vision pricing** — [Claude vision docs](https://platform.claude.com/docs/en/build-with-claude/vision); prices from Anthropic's model table cached 2026-10-06 ([pricing](https://claude.com/pricing)):
  - Image tokens = ⌈w/28⌉ × ⌈h/28⌉; a 1000×1000 image = **1,296 tokens**; capped at 1,568 (standard tier) or 4,784 (high-res tier) tokens.
  - Haiku 5.5: **$0.10 / $0.50** per MTok in/out. Sonnet 5.5: **$2 / $10**. Opus 5.5: **$4 / $20**.
- **Gemini pricing** — [Gemini API pricing](https://ai.google.dev/gemini-api/docs/pricing), paid tier per MTok in/out:
  - gemini-2.5-flash-lite: $0.10 / $0.40.
  - gemini-3.1-flash-lite: $0.25 / $1.50.
  - gemini-2.5-flash: $0.30 / $2.50.
  - gemini-3.6/3.8-flash: $0.75 / $3.75 until 31 Dec 2026, then $1.50 / $7.50 from 1 Jan 2027.
  - The page does not state image token counts for these models.
- **OpenAI**: no official per-image figure retrieved; aggregators conflict (e.g. GPT-5 mini $0.25/$2 per MTok) — [pricepertoken](https://pricepertoken.com/pricing-page/model/openai-gpt-5-image-mini). Treat as unverified.
- **Passio**: 3–5k tokens per image at $2.35–2.50/M ≈ $0.007–0.0125 per image. Its FAQ figure of 20–30k tokens would mean ≈ $0.05–0.075 per image — [Passio](https://www.passio.ai/cost-breakdown).

### Inferences
- **Cost per photo (estimate).** Assumptions:
  - Client downscales the photo to ~1000×1000 → 1,296 image tokens (Claude formula). The same count is assumed for Gemini, which is probably conservative.
  - +1,000 text input tokens (system prompt + schema + short candidate list from Coolinary's catalogue).
  - Output: 1,000 tokens for models that think by default (Claude Haiku/Sonnet/Opus 5.5, Gemini 3.x/2.5 Flash), 300 for Gemini 2.5 Flash-Lite.

  | Model | Arithmetic | $/photo | €/photo |
  |---|---|---|---|
  | Gemini 2.5 Flash-Lite | 2,296 × 0.10/1M + 300 × 0.40/1M | **$0.00035** | €0.0003 |
  | Claude Haiku 5.5 | 2,296 × 0.10/1M + 1,000 × 0.50/1M | **$0.00073** | €0.0006 |
  | Gemini 3.1 Flash-Lite | 2,296 × 0.25/1M + 1,000 × 1.50/1M | **$0.0021** | €0.0018 |
  | Gemini 2.5 Flash | 2,296 × 0.30/1M + 1,000 × 2.50/1M | **$0.0032** | €0.0027 |
  | Claude Sonnet 5.5 | 2,296 × 2/1M + 1,000 × 10/1M | **$0.0146** | €0.0125 |
  | Claude Opus 5.5 | 2,296 × 4/1M + 1,000 × 20/1M | **$0.029** | €0.025 |

- **Per heavy photo-logger per month (estimate):** 3 photos/day × 20 days = 60, +20% retries/corrections = **72 calls**.
  - Haiku 5.5: 72 × $0.00073 = **$0.05 (≈ €0.05)**.
  - Gemini 2.5 Flash: **$0.23 (≈ €0.20)**.
  - Sonnet 5.5: **$1.05 (≈ €0.90)**.
  - At 10k MAU with 30% heavy photo-loggers (3,000 × 72 = 216k photos/month):
    - Haiku ≈ $158 ≈ €136/month.
    - Gemini 2.5 Flash ≈ $689 ≈ €592/month.
    - Sonnet ≈ $3,152 ≈ €2,711/month.
  - At 100k MAU: ×10 (Haiku ≈ €1,355; Gemini Flash ≈ €5,920; Sonnet ≈ €27,100/month).
  - On a ~€3–5/month Premium+AI-style subscription, small models are comfortably affordable. Sonnet-class models are only viable if usage is capped (e.g. N photos/day).
- **Accuracy design for Coolinary.** The evidence (Case 3 vs Case 1) says that giving the model the *expected* ingredients and amounts cuts error by more than half. So:
  - (1) Planned and cooked meals should be logged from the recipe (deterministic, ~0 marginal cost).
  - (2) Photo AI should be a fallback for off-plan meals.
  - (3) The model should output *foods + estimated grams* and let Coolinary's nutrition engine compute kcal from BLS/Ciqual. This keeps the "AI never states calorie numbers" principle in spirit: the AI estimates quantities, the tables state energy.
  - Large-portion under-estimation is the main known bias. Grams should be shown as editable.
- **Voice / text logging** ("2 eggs and a slice of rye bread") is a text-only LLM call, roughly 10× cheaper than a photo. It reuses the same food+grams → nutrition engine pipeline.
- Industry status: AI logging is now a monetisation lever (Fitatu gates it behind Premium+AI; MyFitnessPal bought Cal AI). It is not required for a credible *lite* launch.

### Gaps
- No peer-reviewed accuracy data for Fitatu AI, Yazio AI, MFP Meal Scan, Cal AI or SnapCalorie specifically; vendor claims were not found or not verified.
- No 2026 study on Polish or German dishes (pierogi, bigos, Eintopf), where Western-trained models may err more.
- Gemini per-image token count not confirmed from the pricing page. OpenAI official image pricing not retrieved.

---

## 4. Compliance: GDPR Art. 9, MDR boundary, app stores, AI Act, eating-disorder safeguards, conflicts with Coolinary's rules

### Takeaway
A food diary tied to a person's calorie target, especially with weight, should be treated as **health data under GDPR Art. 9**. The CJEU reads "data concerning health" broadly. That means explicit consent, a DPIA, per-person privacy inside a shared household, and no diary for children.

Lifestyle positioning under MDCG 2019-11 Rev.1 (June 2025) still works if Coolinary makes no disease or treatment claims.

Store-specific requirements:
- Google Play: a Health apps declaration (nutrition/weight management category).
- Apple: data and methodology disclosure for accuracy claims (1.4.1). No health data in iCloud and no ad use (5.1.3).

Two of Coolinary's principles break in a "full clone": "no weight history" and "AI never states calories". A lite design can keep both.

### Cited Findings
- **CJEU C-21/23 *Lindenapotheke*** (4 Oct 2024, Grand Chamber): "data concerning health" must be interpreted broadly. Data that allow conclusions about an identifiable person's health status are health data, even pharmacy *orders* of non-prescription medicines — [Taylor Wessing](https://www.taylorwessing.com/en/insights-and-events/insights/2024/10/ecj-lindenapotheke); [Bird & Bird](https://cm.twobirds.com/en/insights/2024/global/feeling-unwell-after-the-cjeus-lindenapotheke-decision).
  - Germany's BGH applied it in January 2025 (I ZR 222/19) — [Taylor Wessing](https://www.taylorwessing.com/en/insights-and-events/insights/2025/03/bgh-lindenapotheke).
  - The same judgment lets **competitors sue** under unfair-competition law for GDPR breaches. This is relevant in Germany's Abmahnung culture — same sources.
- **GDPR statutory text** — [Regulation (EU) 2016/679](https://eur-lex.europa.eu/eli/reg/2016/679/oj) (statutory text; not re-fetched this session):
  - Art. 9(2)(a): explicit consent as the lawful basis for special-category data.
  - Art. 35(3)(b): DPIA required for large-scale processing of special categories.
  - Art. 8: member-state digital-consent age between 13 and 16.
- **MDCG 2019-11 Rev.1** published 17 June 2025 — [European Commission](https://health.ec.europa.eu/latest-updates/update-mdcg-2019-11-rev1-qualification-and-classification-software-regulation-eu-2017745-and-2025-06-17_en).
  - It clarifies modular software and borderline cases. Qualification still hinges on **intended purpose** — [Emergo by UL](https://www.emergobyul.com/news/european-revision-primary-software-guidance-mdcg-2019-11-revision-1-small-changes-meaningful); [OpenRegulatory](https://openregulatory.com/mdcg/mdcg-2019-11).
  - None of the search results addressed lifestyle or well-being apps directly.
- **Google Play.**
  - All apps must complete the **Health apps declaration**. Its categories include nutrition and weight management (dietary intake tracking, meal planning, weight-loss goals) — [Google Play Help](https://support.google.com/googleplay/android-developer/answer/14738291).
  - A policy preview says Google will incorporate MDCG guidance for medical-device apps targeting the EU — [Play policy preview](https://support.google.com/googleplay/android-developer/answer/13996823) (snippet; no effective date seen).
  - Health Connect access requires a declaration — [Android developers](https://developer.android.com/health-and-fitness/health-connect/publish).
- **Apple App Review Guidelines** — [Apple guidelines](https://developer.apple.com/app-store/review/guidelines/):
  - **1.4.1**: "Apps must clearly disclose data and methodology to support accuracy claims relating to health measurements… if the level of accuracy or methodology cannot be validated, we will reject your app"; apps "should remind users to check with a doctor".
  - **5.1.3(i)**: health/fitness data (incl. HealthKit) may not be used for advertising or data mining.
  - **5.1.3(ii)**: apps "may not store personal health information in iCloud".
  - Developer-forum reports show nutrition and weight-loss apps rejected under 1.4.1 for **missing citations** for health recommendations — [Apple forum 767340](https://developer.apple.com/forums/thread/767340); [Apple forum 696430](https://developer.apple.com/forums/thread/696430) (older thread, 2021–22; possibly outdated).
- **EU AI Act Art. 50** transparency obligations (e.g. telling users they are interacting with an AI system) apply from **2 Aug 2026**, regardless of risk class. Marking of AI-generated content has extended deadlines for systems already on the market — [Cooley](https://www.cooley.com/news/insight/2026/2026-08-03-eu-ai-act-transparency-obligations-take-effect-2-august-2026); [Sidley](https://datamatters.sidley.com/2026/06/24/eu-ai-act-transparency-obligations-preparing-for-compliance-by-2-august-2026/). Sources disagree on exactly which Art. 50 sub-paragraphs apply on that date.
- **Eating-disorder risk.**
  - Evidence that calorie trackers are associated with ED symptoms is **correlational**. Some cohorts self-report harm (e.g. a 2017 sample where 73% of ED patients using MyFitnessPal said it contributed). An RCT in undergraduate women without prior symptoms found no difference — [PMC8485346](https://pmc.ncbi.nlm.nih.gov/articles/PMC8485346); [National Center for Health Research](https://www.center4research.org/fitness-tracking-apps-eating-disorders/). The 2017 data is pre-2024 and possibly outdated.
  - Recommended safeguards: screening at sign-up and check-ins, signposting support services, blocking unhealthy goals — same sources. No independent evaluation of safeguard effectiveness was found.

### Inferences
- **GDPR classification (inference).** Under *Lindenapotheke*'s broad reading, a per-person intake log next to a computed calorie target (from age, sex, height, weight) lets one infer weight status, dieting behaviour and possibly ED patterns. Treat it as Art. 9 data. Even the lite "eaten" flag on planned meals becomes intake logging once tied to a named household member.

  Required:
  - Explicit, granular, separately revocable consent for the diary.
  - A DPIA before launch.
  - Data minimisation (local-first storage where possible).
  - An Art. 30 processing record.
  - Processor DPAs with the AI vendor, including EU/EEA data residency or SCCs.
  - Export and delete.
  - Competitor-lawsuit risk in DE makes the privacy policy and consent UX worth a lawyer review (≈ €1–3k one-off; unsourced estimate).
- **Household privacy.** Coolinary is household-shared. Diaries must be **private per adult by default**: no "Mum can see Dad's intake". Each adult consents for themselves.
- **Children.** Germany and Poland are understood to use 16 as the digital-consent age (not verified this session). Given Coolinary's special non-diet rules for children, the safest rule is **no diary, no intake tallies and no calorie numbers on child profiles** (ideally under 18). Children's meals still count toward the household plan and stock.
- **MDR (inference).** A diary for healthy adults managing weight as a lifestyle goal stays outside MDR if marketing and UX avoid claims to treat, prevent or monitor disease. Examples that would push it in: "for obesity", "for diabetics", carb counting for insulin dosing, PCOS or IBS protocols. Coolinary's MDCG 2019-11 lifestyle positioning carries over. The diary copy must never claim clinical accuracy (also needed for Apple 1.4.1).
- **DiGA (Germany).** Not relevant. DiGA is a reimbursement pathway for CE-marked medical-device apps, which a lifestyle app deliberately is not (general regulatory knowledge, not re-verified this session).
- **AI Act (inference).** Food recognition for self-logging is not an Annex III high-risk use, so it is minimal or limited risk. The Art. 50 duty to disclose AI interaction applies to Coolinary's AI meal composer and to any AI logging assistant: label AI-generated estimates clearly.
- **Apple 5.1.3(ii) and "weight kept on the phone".** Keeping weight local is compatible, but the app must not sync weight or diary to iCloud (CloudKit / iCloud key-value store). Device-backup nuances are worth a check.
- **Conflicts with Coolinary's principles, and resolutions:**

  | Coolinary rule | Full clone needs | Lite-compatible resolution |
  |---|---|---|
  | No food diary | Persistent per-day intake log | Unavoidable: any tracker *is* a diary. Make it opt-in, per-adult, private, deletable, with short default retention (e.g. 90 days) |
  | No weight history | Weight log and chart (table stakes) | Delegate to **Apple Health / Health Connect**: read the latest weight for targets, optionally show the OS chart. Coolinary stores only the current value on-device |
  | AI never states calorie numbers | AI photo logging outputs kcal | AI outputs foods + grams only; kcal computed deterministically from BLS/Ciqual and shown as table-derived values |
  | Child-safety, non-diet rules | Diary for everyone | Hard-disable diary and tallies on child profiles |
  | Lifestyle (MDCG) positioning | Competitors market "lose X kg" | Avoid disease claims; add Apple-required "consult a doctor" note and methodology/citation page (Mifflin–St Jeor, BLS/Ciqual sources) |

- **ED safeguards for launch (synthesis).**
  - Keep the 1,200 kcal adult floor.
  - Block goal weights that imply BMI < 18.5 (if weight goals are added at all).
  - Cap the weekly loss rate.
  - "Hide numbers" mode (Coolinary's existing default can be numbers-hidden).
  - No streak-shaming.
  - Optional screening question at diary opt-in, with signposting:
    - PL: Ogólnopolski Telefon Zaufania dla Dzieci i Młodzieży 116 111 (helpline; number not verified this session).
    - DE: BZgA eating-disorder helpline (not verified).
  - No "earn back" calories from exercise by default.

### Gaps
- No official EDPB statement specifically on fitness or diet apps was found this session (older WP29 guidance on health data predates 2024 and was not retrieved).
- Apple does not publish a minimum-calorie rule; I found no evidence of one. The "1,200 kcal floor" is industry practice and Coolinary's own rule, not a store requirement.
- The Play policy update incorporating MDCG has no confirmed effective date.
- The exact current PL/DE Art. 8 digital-consent ages were not re-verified.

---

## 5. Engineering effort for a solo developer (cross-platform stack): lite vs full

### Takeaway
On top of Coolinary's existing nutrition engine, barcode, catalogue and targets:
- A **lite** tracker ("eat the plan" + quick-add + snacks + barcode-to-diary + daily summary + consent and privacy) is about **9–13 person-weeks**.
- A **"de facto Fitatu/Yazio clone"** (full search, custom foods, servings, history and charts, weight via HealthKit/Health Connect, wearables, reminders, AI photo/voice, moderation tooling) adds about **24–37 person-weeks**, for **33–50 person-weeks total** (8–12 months solo).

Public build stories match this:
- MacroFactor started with two engineers on Flutter.
- Cal AI launched with three co-founders.
- Solo "MVP in ~3 months" stories describe much thinner apps.

### Cited Findings
- MacroFactor: built with Flutter + Firebase to ship one codebase to both stores. A co-founder said they started with **two engineers** — [Google Cloud case study](https://cloud.google.com/customers/macrofactor) (snippet). Launched 2021; bootstrapped — [MacroFactor team](https://macrofactor.com/team/).
- Cal AI:
  - Launched May 2024 by three co-founders (self-reported) — [Latka](https://getlatka.com/companies/calai.app).
  - About $2M revenue reported by March 2025 — [YPulse](https://www.ypulse.com/newsfeed/2025/03/20/two-high-schoolers-founded-cal-ai-a-calorie-tracking-app-with-millions-of-downloads-and-2m-in-revenue/).
  - 7 employees plus contractors at acquisition — [TechCrunch](https://techcrunch.com/2026/03/02/myfitnesspal-has-acquired-cal-ai-the-viral-calorie-app-built-by-teens).
  - Latka separately claims 17 people and $35M (inconsistent, self-reported).
- Solo examples (anecdotal):
  - "NutriTrack" was built to ship solo in about 3 months — [Contra case study](https://contra.com/p/evrl2omW-nutri-track-calorie-tracking-app-development).
  - Other indie calorie apps report $1.5k MRR two months after launch, or years with no traction — [search summary sources, e.g. Indie Hackers interview](https://www.indiehackers.com/interview/body-tracker-app-0d0c5cfd05) (snippet; unverified).

### Inferences
Person-week (pw) estimates assume an experienced Flutter or React Native developer, an existing backend with auth and sync, and testing on both platforms. Ranges reflect unknowns about the current codebase.

**(a) Lite on top of the existing plan:**

| # | Item | New / existing | Effort |
|---|---|---|---|
| L1 | Intake data model (person, date, slot, source = planned dish / product / quick-add / custom; servings; nutrient snapshot), daily aggregation | New | 1–1.5 pw |
| L2 | "Mark planned meal eaten" + portion multiplier (½, 1, 1½, custom g); hook into existing "cooked it" (shared dish split by person portions) | New on existing recipes/portions | 1–1.5 pw |
| L3 | Quick-add kcal (+ optional P/F/C) and snack slot | New | 0.5 pw |
| L4 | Barcode/catalogue → diary with serving picker (per 100 g / per package / per piece) | Scanner and lookups exist; serving UI new | 1–1.5 pw |
| L5 | Daily summary vs per-person target, edit/delete, recents, "copy yesterday" | Targets exist; UI new | 1.5–2 pw |
| L6 | Consent (Art. 9 explicit), per-adult privacy in household, export/delete, retention, DPIA + privacy policy + Art. 30 record | New | 1.5–2 pw |
| L7 | Child-profile hard-disable, ED safeguards, "hide numbers" mode, methodology/citations page (Apple 1.4.1) | Partly existing rules | 0.5–1 pw |
| L8 | Offline-first entry queue + sync conflict handling | Depends on existing sync | 0.5–1.5 pw |
| L9 | QA, store submission (Play Health apps declaration update, App Privacy labels), analytics | New | 1–1.5 pw |
| | **Lite total** | | **≈ 9–13 pw (≈ 2–3 months solo)** |

**(b) Full clone (additional to lite):**

| # | Item | Effort |
|---|---|---|
| F1 | Full-text food search across OFF + BLS/Ciqual + crowd + custom; PL/DE stemming, typo tolerance, ranking, dedupe; self-hosted Typesense/Meilisearch; OFF delta import pipeline | 3–5 pw |
| F2 | User custom foods (label entry incl. photo of nutrition panel) + crowd submission | 1.5–2 pw |
| F3 | Household measures per food (łyżka, szklanka, Scheibe, Esslöffel) and unit conversions | 1–2 pw |
| F4 | History: calendar, weekly/monthly charts, macro/micro detail, averages | 2–3 pw |
| F5 | Weight via HealthKit / Health Connect (read/write, OS-owned history); optional in-app chart reading OS data | 1.5–2 pw |
| F6 | Wearables: active energy/steps import (HealthKit, Health Connect), write nutrition back; permission UX; Health Connect declaration | 2–3 pw (direct Garmin/Fitbit/Withings APIs +1.5–2 pw each) |
| F7 | Water tracker | 0.5 pw |
| F8 | Reminders, notifications, streaks (non-shaming) | 1–1.5 pw |
| F9 | AI photo + text/voice logging: capture, LLM call with structured output (foods + grams), catalogue matching, editable result UI, rate limiting/paywall gating, eval set of ~100 PL/DE meals | 4–6 pw |
| F10 | Moderation/admin tooling for crowd products (review queue, duplicate merge, reports, ODbL separation) | 2–3 pw |
| F11 | Localisation PL/DE/EN of all new UI, accessibility, hardening, beta cycle | 2–3 pw |
| F12 | Premium paywall for tracker features (if not already built); must use store IAP correctly (cf. Cal AI removal) | 1–2 pw |
| | **Full-clone additional** | **≈ 22–33 pw; with integration overhead ≈ 24–37 pw** |
| | **Full total (lite + full)** | **≈ 33–50 pw (≈ 8–12 months solo)** |

- A solo developer cannot realistically match the Fitatu or Yazio *database*, whatever the engineering effort. That is an operations and data problem (see section 6), not a coding problem.

### Gaps
- No primary, dated build log of a solo developer who shipped a full PL/DE tracker was found. Indie stories are anecdotal and unverified.
- MacroFactor's launch-era headcount and build duration are not documented beyond "started with two engineers".

---

## 6. Ongoing costs: database curation, user-submitted products, support, dietitian review

### Takeaway
Run costs are dominated by **human time on data quality and support**, not infrastructure. With OFF + crowd data and small-model AI, cash costs stay low:
- 10k MAU: roughly **€50–300/month**.
- 100k MAU: roughly **€500–2,500/month**.

Founder time for moderation and support grows from about **10–20 h/month at 10k MAU** to about **80–170 h/month at 100k MAU**, i.e. part-time to full-time help.

### Cited Findings
- OFF data is free under ODbL with attribution and share-alike — [OFF terms](https://world.openfoodfacts.org/terms-of-use); [ODbL](https://opendatacommons.org/licenses/odbl/1-0/).
- Commercial DB prices:
  - Edamam: $14 / $69 / $299 per month — [Edamam](https://developer.edamam.com/food-database-api).
  - FatSecret Premier: quote-only; Premier Free available for small start-ups — [FatSecret](https://platform.fatsecret.com/api-editions).
  - Passio: about $0.625 per average active user/month at $2.50/M tokens — [Passio](https://www.passio.ai/cost-breakdown).
- AI per-photo prices — [Claude vision](https://platform.claude.com/docs/en/build-with-claude/vision); [Gemini pricing](https://ai.google.dev/gemini-api/docs/pricing). Computed in section 3.

### Inferences
All figures below are estimates; assumptions stated.

- **User-submitted product verification.** Assume 1% of MAU submit 1 new product/month at ~2 min review each (label photo vs entered values; auto-checks such as kcal ≈ 4P + 4C + 9F ±15% flag outliers).
  - 10k MAU: 100 × 2 min ≈ **3.3 h/month**.
  - 100k MAU: 1,000 × 2 min ≈ **33 h/month**.
  - In Poland, early submission rates will be higher because OFF coverage is thin: plan 2–3× that for the first 6 months.
  - Option: an OCR + LLM auto-extract from the label photo using a small vision model (≈ €0.001/label per section 3), with a human only on flagged items.
- **Data error reports / duplicates.** About +50% of the above.
- **Support.** Assume 1–2% of tracker-active MAU contact support monthly at ~5 min/ticket.
  - 10k MAU: 100–200 tickets ≈ **8–17 h/month**.
  - 100k MAU: ≈ **80–170 h/month**.
  - Trackers produce more "wrong kcal / missing product / sync" tickets than meal planners.
- **Dietitian review.** One-off review of targets, floors, ED safeguards, copy and the methodology page: ~10–20 h. At an assumed PL/DE freelance dietitian rate of €40–80/h (unsourced), that is **€400–1,600 one-off**, plus ~2–4 h/quarter for content updates.
- **Infra.**
  - Diary rows: 10k MAU × 30% loggers × 5 entries/day × 30 days ≈ 450k rows/month. Trivial for Postgres.
  - Search cluster + OFF mirror: ≈ €20–100/month at 10k MAU and €100–400/month at 100k MAU (unsourced hosting estimate).
- **AI logging** (section 3, small model, 30% heavy users): ≈ €65–600/month at 10k MAU; ≈ €650–6,000/month at 100k MAU. Gate it behind premium to keep gross margin positive.
- **Legal.** DPIA + privacy policy + consent UX review ≈ €1–3k one-off (estimate; Lindenapotheke competitor-lawsuit exposure in DE).

### Gaps
- No public data on actual moderation or support ratios at Fitatu, Yazio or OFF. The ratios above are assumptions to be replaced by beta telemetry.
- No sourced PL/DE dietitian hourly rates.

---

## 7. Launch checklist: what Coolinary already has vs what is new; lite vs full clone

### Takeaway
**Lite ("eat the plan")** is about 2–3 months solo and well under €300/month at 10k MAU. It reuses most of Coolinary's plan and preserves its principles if the diary is opt-in, per-adult and private, with weight delegated to the OS.

**Full clone** is about 8–12 months solo. It forces reversing "no weight history" (or delegating it), relaxing the "AI never states numbers" rule, and running a crowd-curation operation that cannot close the Polish database gap with Fitatu within a year.

### Cited Findings
- Database gap figures: OFF PL 38,498 total / 26,381 nutrition-complete vs Fitatu PL >500k — [OFF API](https://world.openfoodfacts.org/api/v2/search?countries_tags=en:poland&page_size=1&fields=code); [Fitatu help](https://www.fitatu.com/pl/help/czy-fitatu-posiada-najwieksza-baze-polskich-produktow-i-potraw).
- AI logging accuracy improves sharply when the model is given the ingredient list (53 vs 123 kcal MAE) — [Nutrients 2025](https://pmc.ncbi.nlm.nih.gov/articles/PMC12655113/).
- Competitor monetisation of AI (Fitatu Premium+AI 49.99–64.99 zł/month; conflicting reports) — [Telepolis](https://www.telepolis.pl/tech/aplikacje/fitatu-ai-kalorie-ze-zdjecia-posilku); [Tabletowo](https://www.tabletowo.pl/?p=960186).

### Inferences
**Checklist.** "Have" = in Coolinary's plan already. Effort is in person-weeks. Cash cost is per month at 10k / 100k MAU unless one-off.

| Item | Have? | Lite | Full | Effort | Cash cost |
|---|---|---|---|---|---|
| Nutrition engine (Ciqual/BLS/USDA), per-portion kcal and macros | **Have** | ✓ | ✓ | 0 | 0 |
| Per-person targets (Mifflin–St Jeor, 1,200 floor, child rules) | **Have** | ✓ | ✓ | 0 (macro split UI +0.5) | 0 |
| Barcode + own/crowd/OFF/GS1 lookups | **Have** | ✓ | ✓ | 0 | OFF free |
| "Cooked it" stock decrement | **Have** | ✓ (log hook) | ✓ | incl. L2 | 0 |
| Intake model + daily aggregation | New | ✓ | ✓ | 1–1.5 | ~0 |
| Mark planned meal eaten + portion multiplier | New | ✓ | ✓ | 1–1.5 | 0 |
| Quick-add kcal + snacks slot | New | ✓ | ✓ | 0.5 | 0 |
| Barcode → diary serving picker | New (UI) | ✓ | ✓ | 1–1.5 | 0 |
| Daily summary, edit, recents, copy yesterday | New | ✓ | ✓ | 1.5–2 | 0 |
| Art. 9 consent, per-adult privacy, export/delete, DPIA, policy | New | ✓ | ✓ | 1.5–2 | €1–3k one-off legal |
| Child hard-disable, ED safeguards, hide-numbers, methodology page | Partly | ✓ | ✓ | 0.5–1 | €0.4–1.6k one-off dietitian |
| Offline queue + sync | Depends | ✓ | ✓ | 0.5–1.5 | 0 |
| QA + store declarations | New | ✓ | ✓ | 1–1.5 | 0 |
| Full food search infra (OFF mirror + BLS/Ciqual + crowd) | Partial | — | ✓ | 3–5 | €20–100 / €100–400 |
| Custom foods + crowd submission | Partial (crowd catalogue) | — (optional) | ✓ | 1.5–2 | moderation time |
| Household measures / units | New | — | ✓ | 1–2 | 0 |
| History and charts | New | — (7-day mini view optional) | ✓ | 2–3 | 0 |
| Weight via HealthKit / Health Connect | New (conflict) | — | ✓ | 1.5–2 | 0 |
| Wearables (activity kcal) | New | — | ✓ | 2–3 | 0 |
| Water | New | — | ✓ | 0.5 | 0 |
| Reminders / streaks | New | optional (1 daily reminder) | ✓ | 1–1.5 | 0 |
| AI photo/voice logging (foods + grams → engine) | New (conflict) | — (or premium beta) | ✓ | 4–6 | €65–600 / €650–6,000 (small model, 30% heavy users) |
| Moderation tooling | New | — | ✓ | 2–3 | 3–33 h/month founder time |
| Localisation + hardening | Partial | incl. | ✓ | 2–3 | 0 |
| Paywall via store IAP | ? | optional | ✓ | 1–2 | store fees |
| Commercial DB (FatSecret PL/DE or Edamam) | — | — | optional | 1 (integration) | Edamam €59–257 at 10k; FatSecret quote |
| **Totals** | | **≈ 9–13 pw** | **≈ 33–50 pw** | | Lite cash ≈ €20–150/month at 10k; Full ≈ €100–1,000/month at 10k excl. salaries |

**Recommended minimum to launch a *credible* lite tracker:**
- L1–L9, plus one daily reminder.
- A 7-day mini history.
- "Add product" through the existing crowd catalogue.

Market it as "zjedzone z planu" / "aus dem Plan gegessen" (eaten from the plan) tracking, not as a Fitatu replacement.

**Upgrade triggers to full:**
- Over 30% of logged entries are off-plan.
- Search-miss rate over 15% in Poland.
- Weight-feature requests in the top 3 of feedback.

**Germany-first** is the lower-risk market for any database-dependent features (OFF DE ≈ 11× OFF PL).

### Gaps
- Codebase specifics (existing sync, paywall, search) are unknown, so the effort ranges are wide.
- Conversion and retention uplift from adding a diary to a meal planner is not quantified here; that belongs to the ROI workstream.
