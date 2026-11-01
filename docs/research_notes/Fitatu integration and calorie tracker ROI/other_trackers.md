# Integration routes into calorie trackers other than Fitatu (Yazio, MyFitnessPal, Lifesum, Cronometer, FatSecret, Lose It!, MacroFactor, Samsung Health, Cal AI)

Research date: 30 Oct 2026. Pages behind 403 (Yazio help centre, Lose It Zendesk, FatSecret app help, MFP support, Yazio feedback board) could not be fetched directly. For those, the notes use search-engine snippets and label them that way. App Store prices and ranks were read on 30 Oct 2026 from the PL and DE storefronts.

## 1. Public or partner API: status in 2026, access requirements, practicality

### Takeaway
Only **FatSecret** offers a self-serve, documented API that can **write to a user's food diary**: OAuth 1.0 3-legged "profile" methods such as `food_entry.create`, with free tiers for small start-ups. **Lose It!** has an official partner API, but it is beta, **read-only**, and limited to businesses. **Cronometer** shares its API only with a few selected B2B partners. **MyFitnessPal** says on its API site that it is not accepting API requests. **Yazio, Lifesum and MacroFactor** have no public or partner API. **Samsung Health** has an official SDK that can write Nutrition, but only after Samsung approves a partner request.

### Cited Findings

**MyFitnessPal (MFP)**
- The official API site says: "We are not accepting requests for API access at this time." It gives no criteria or application process. The footer has a "Licensing" contact, partners@myfitnesspal.com. — [myfitnesspalapi.com](https://www.myfitnesspalapi.com/)
- Background (2012, outdated): MFP launched a partner API in Oct 2012 with Fitbit, Withings, Runtastic, Endomondo and others. — [TechCrunch 2012](https://techcrunch.com/2012/10/16/myfitnesspal-api)
- A third-party API catalogue describes OAuth 2.0 partner credentials with diary and measurement scopes, and says the program is "closed to new applicants". — [apis.io](https://apis.io/security/myfitnesspal/myfitnesspal-authentication/) (secondary source)
- Vital (tryvital.io) lists a MyFitnessPal connector. This suggests a few aggregators still have legacy partner access. — [tryvital.io](https://tryvital.io/wearables-api/myfitnesspal) (vendor page, not verified)
- MFP builds integrations top-down with large platforms. On 7 Jan 2026 it launched an app inside ChatGPT (ChatGPT Health) with insights, data and recipes. Food logging from ChatGPT was not available at launch and was described only as a possible future feature. — [Healthcare Technology Report](https://thehealthcaretechnologyreport.com/myfitnesspal-integrates-nutrition-app-into-new-chatgpt-health-experience/); [press release mirror](https://business.mammothtimes.com/mammothtimes/article/gnwcq-2026-1-7-myfitnesspal-launches-app-in-new-chatgpt-health-experience)
- The Terms of Service (effective 18 Mar 2025) require a "separate agreement with us … in advance" for any commercial use "other than through the Commercial Tools". "APIs" for website owners are named as one example of MFP-branded tools. — [MFP Terms](https://www.myfitnesspal.com/terms-of-service)

**FatSecret Platform API**
- Editions — [platform.fatsecret.com/api-editions](https://platform.fatsecret.com/api-editions):
  - **Basic:** free self sign-up, 5,000 calls/day, US dataset only, English, attribution required.
  - **Premier Free:** free after verification for start-ups, non-profits and students. Unlimited calls, US dataset; non-US data is "available at a discount for eligible users". Attribution required.
  - **Premier:** price "upon request, based on market access". More than 62 countries, 26 languages, no attribution, SLA, and the right to request database edits.
- Premier Free start-up eligibility: "a company earning less than US$1 million in annual revenue" that has raised less than US$1M in funding. Add-ons are discounted 50% for start-ups and non-profits. — [api-editions](https://platform.fatsecret.com/api-editions)
- All editions list: food search, autocomplete, recipe database, recipe/meal builder, custom food and exercise tracker, **food/exercise diaries** and weight tracking. — [api-editions](https://platform.fatsecret.com/api-editions)
- Image recognition and NLP are optional add-ons billed per use, and are "not included in any of our standard API tiers (e.g. Basic, Premier Free or Premier)" (FatSecret forum reply, March 2026). — [Google Group](https://groups.google.com/g/fatsecret-platform-api/c/B4SziL4FpCg)
- Localization is "a premium feature only made available to select accounts". Region **PL** (default language `pl`) and **DE** (default `de`) are both listed, and Polish and German are both among the 27 language codes. — [Localization guide](https://platform.fatsecret.com/docs/guides/localization)
- **Writing to a diary:** "Profile – Food Diary: Create Entry By Id (v1)" is `POST https://platform.fatsecret.com/rest/food-entries/v1` (legacy name `food_entry.create`). Required parameters: `food_id`, `food_entry_name`, `serving_id`, `number_of_units`, `meal` (breakfast/lunch/dinner/other) and `date`. It sits under "Profiles – 3-Legged OAuth" with an OAuth 1.0 badge, and it is **not** marked Premier-exclusive. — [food_entry.create docs](https://platform.fatsecret.com/docs/v1/food_entry.create)
- Profiles are how the platform stores per-user data. `profile.create` and `profile.get_auth` issue per-user OAuth tokens. Methods like `food_entry.create` act on a user's Food Diary. — [Authentication guide](https://platform.fatsecret.com/docs/guides/authentication)
- The pyfatsecret wrapper docs say full 3-legged OAuth "can attain an access token for a profile linked to a fatsecret.com user account". — [pyfatsecret usage](https://pyfatsecret.readthedocs.io/en/latest/usage.html) (third-party docs)
- FatSecret Brand Tools let brands "Manage and Publish Nutrition Data". — [platform.fatsecret.com/brands](https://platform.fatsecret.com/brands)

**Lose It!**
- The official partner page describes a beta API. Food logs (calories and macros) are among the data types. It uses OAuth 2.0 and webhooks that notify partner servers when user data is available or updated. The page says "the API is read-only and we do not currently support write access." Docs are at developers.loseit.com. — [Lose It! API page](https://d2qjq4b9ehc3ag.cloudfront.net/api)
- Access is requested through a survey. Lose It! is "only considering API partners for business/commercial applications" and does not grant access to individuals who want their own logging data. — [Lose It! API survey](https://www.surveymonkey.com/r/ZGYMDC7); [second form](https://www.surveymonkey.com/r/PGWGPS3)
- developers.loseit.com could not be fetched ("unable to fetch"), so pricing and geographic limits are unverified.

**Cronometer**
- There is no public self-serve API. The Cronometer Pro FAQ says "Everfit, Kalix & Practice Better have access to Cronometer's API". — [Cronometer Pro FAQ](https://cronometer.com/blog/cronometer-pro-faq)
- Users manage a "Connected Apps" list of API partners they have granted access to. — [Cronometer support](https://support.cronometer.com/hc/en-us/articles/360018867471)
- A long-running forum thread asks for a public API. A user there claims API access was added "only for enterprise level plans" (not officially confirmed). — [Cronometer forum](https://forums.cronometer.com/discussion/3365/any-chance-a-public-api-for-cronometer-will-be-available)
- G2 describes Cronometer Pro Plus (institutional tier) as including API access. — [G2](https://www.g2.com/products/cronometer-pro/discuss) (secondary)

**Yazio**
- There is no public API. A "Expose public API" request on Yazio's own feature board has about 424 votes and the status "Under consideration". A separate request asks for a public API or MCP server. — [Yazio feature board](https://yazioen.featureupvote.com/suggestions/28294/expose-public-api); [feedback-en.yazio.com](https://feedback-en.yazio.com/suggestions/713173/a-public-api-andor-mcp-server-integration) (search snippets; the pages returned 403)

**Lifesum**
- An API directory lists no public developer program, documentation or partner portal. — [API Evangelist](https://providers.apievangelist.com/providers/lifesum/) (secondary)

**MacroFactor**
- A MacroFactor developer wrote on r/MacroFactor (Jan 2023, possibly outdated) that an API was "not part of our near, medium or long term plans at this moment, although I won't say never". — [Reddit mirror](https://axebps-redlib.hf.space/r/MacroFactor/comments/10b7as9/public_facing_api/j4ba48f/?context=3)

**Samsung Health (food log)**
- The Samsung Health Data SDK supports insert, update and delete for **Nutrition** and water intake, among other types. — [Release notes](https://developer.samsung.com/health/data/release-note.html); [NutritionType API](https://developer.samsung.com/health/data/api-reference/-shd/com.samsung.android.sdk.health.data.request/-data-type/-nutrition-type/index.html)
- The FAQ says "To write health data into the Samsung Health, request partnership". After approval, a write access code is issued, "available for one application". Apps that ship to users and read data must also "register your app by submitting a partner request". Developer mode allows testing without registration. — [Samsung Health Data SDK FAQ](https://developer.samsung.com/health/data/faq.html)
- Samsung Health and Health Connect sync nutrition on the phone. A Nutrition item appears in Health Connect when Samsung Health has shared nutrition. — [Samsung dev blog](https://developer.samsung.com/health/blog/en/accessing-samsung-health-data-through-health-connect)

**Cal AI (AI photo logging)**
- MyFitnessPal acquired Cal AI. The article is dated 3 Mar 2026 and the terms were not disclosed. Cal AI had "more than 15 million downloads" and "continues to operate as a separate product within the MyFitnessPal portfolio". The article mentions no partner API. — [The Next Web](https://thenextweb.com/news/myfitnesspal-acquires-cal-ai-the-viral-calorie-tracking-app-built-by-teens)
- A Cal AI App Store listing says it syncs "across iPhone, Apple Watch, and Apple Health (HealthKit)". — [App Store](https://apps.apple.com/us/app/-/id6755081229)

### Inferences
- In practice there are two official ways to push meals: FatSecret's diary API and Samsung Health's Nutrition write (with partner approval). For everything else, the realistic "push" route is writing nutrition to HealthKit or Health Connect, which only some trackers read (see section 5).
- FatSecret Premier Free fits a pre-revenue, bootstrapped EU start-up. The catch is that PL/DE localisation is a Premier or "select accounts" feature, so a request for discounted non-US access is needed.
- `food_entry.create` needs a FatSecret `food_id` and `serving_id`. Coolinary would have to map each dish (or each ingredient) to FatSecret foods, or first create a custom food or recipe through the profile methods. I did not verify whether a custom-food creation method exists.
- Whether entries created on a platform "profile" linked by 3-legged OAuth appear in the user's consumer fatsecret app diary is not stated in the official docs. The pyfatsecret wording suggests they do. Test this before relying on it.
- Lose It!'s read-only API would let Coolinary read what users ate (for example, to adjust targets) but not push meals. Lose It! is US-centric, so its value in PL/DE is likely low (see Gaps).
- MyFitnessPal's 2026 moves (ChatGPT Health app, Cal AI acquisition) show partnerships with large platforms. That is not a signal of a reopened programme for small start-ups.

### Gaps
- MFP: no 2026 statement on criteria for partners@myfitnesspal.com. Whether the Vital/Sahha-style aggregator access still works is unverified.
- Lose It!: pricing, the full data-type list and EU availability (developers.loseit.com was blocked).
- Cronometer: the application process and price for API partner or enterprise access.
- Samsung: whether partner requests are currently accepted, and the eligibility criteria (company size, region). The FAQ does not say.
- Cal AI: no developer page or partner API found. Unclear whether it writes dietary energy to HealthKit or Health Connect, or whether it reads it.

## 2. Recipe import by URL (could Coolinary publish schema.org recipe pages that users import?)

### Takeaway
URL recipe import is documented officially for **MyFitnessPal** (English-language website only), **Cronometer** (Gold-only) and **MacroFactor**. **Yazio** does not have it; it is an open feature request. No evidence was found that **Lifesum** or **FatSecret** have it, and evidence for **Lose It!** conflicts. None of the official help pages say they parse schema.org/Recipe JSON-LD or the `nutrition` property, and none mention Polish or German pages.

### Cited Findings
- **MyFitnessPal:** "With this feature, you can quickly import recipes from a link or add ingredients by hand to match them with the food database. It is currently available to visitors using the English-language version of the website", while other language versions still use the older tool. — [MFP support: Recipe Importer](https://support.myfitnesspal.com/hc/en-us/articles/360032271592-How-does-the-Recipe-Importer-on-the-website-work) (search snippet; page 403)
- MFP has a separate help page for the original recipe tool on non-English sites. — [MFP support: non-English recipe tool](https://support.myfitnesspal.com/hc/en-us/articles/360032271452-How-do-I-create-a-recipe-Original-Recipe-Tool-Non-English-versions-of-the-site)
- MFP Premium gating is disputed. NutriScan says MFP's URL import is free and works in the website and iOS/Android apps. A competitor (Nutrola) says on one page there is no free URL import and on another that it is Premium. — [NutriScan](https://nutriscan.app/blog/posts/best-app-recipe-import-2026-fb4fc6a4f6); [Nutrola](https://nutrola.app/en/blog/free-calorie-tracker-with-recipe-import-2026) (both biased secondary sources)
- **Cronometer** (blog post, 25 Jan 2024):
  - "The Recipe Importer is an exclusive feature available to Gold subscribers."
  - It works on mobile (Foods > Import Recipe) and web (Foods > Custom Recipes > Import Recipe).
  - Ingredients are matched to sources such as USDA and NCCDB. The matching algorithm was retrained and reaches a "92% success rate".
  - Unsupported sites fall back to pasting the recipe text.
  - Structured data and language are not mentioned.
  — [Cronometer blog](https://cronometer.com/blog/recipe-importer/); [Mobile Recipe Importer help](https://support.cronometer.com/hc/en-us/articles/4403966087828-Mobile-Recipe-Importer)
- Cronometer forum users report that some sites fail to import. — [Cronometer forum](https://forums.cronometer.com/discussion/comment/19538)
- **MacroFactor:** "Import from Link" in custom recipes imports "the name, serving quantity, ingredients, and preparation method". The user can then enter the cooked total weight so the recipe can be logged by mass. — [MacroFactor help](https://help.macrofactorapp.com/en/articles/259-import-recipes-from-link)
- **Yazio:** "Import recipes from web by entering the url" is a feature request with about 85 votes, citing MFP as already offering it. — [Yazio feature board](https://yazioen.featureupvote.com/suggestions/90212/import-recipes-from-web-by-entering-the-url) (search snippet)
- Yazio recipes are created manually: Diary > + > ⋯ > "Neues Rezept", with at least 2 foods. — [GIGA](https://www.giga.de/tipp/yazio-eigene-rezepte-erstellen/)
- **Lifesum:** the DE App Store listing mentions only a "große Rezeptbibliothek mit smarten Filtern" (in-house recipe library) and no import. — [App Store DE](https://apps.apple.com/de/app/id286906691)
- **Lose It!:** a low-authority site claims Premium members can paste a recipe link. Competitor blogs say Lose It! has no or limited URL import. The official "Using Recipes in Lose It" article (search snippet) only describes creating a recipe from a logged meal. — [ntro.it.com](https://ntro.it.com/how-do-i-enter-a-recipe-in-the-lose-it-app); [Lose It help](https://loseit.zendesk.com/hc/en-us/articles/47345201948052-Using-Recipes-in-Lose-It)
- **FatSecret app:** no evidence of URL recipe import found. — (searches returned only Apple Health and workout-import help pages)

### Inferences
- Recipe URL importers in general rely on schema.org/Recipe (`recipeIngredient`) markup, as the recipe-scraper ecosystem does. This is general industry knowledge, **not confirmed in any tracker's official docs**.
- None of the trackers say they use the `nutrition` property. They all re-derive nutrition by matching ingredient lines to their own database. Coolinary's calorie numbers would therefore not carry over, and the tracker's totals could differ from Coolinary's.
- The MFP importer is English-site only, and Cronometer matches against USDA/NCCDB. Polish or German ingredient lines ("200 g twarogu półtłustego", "1 EL Rapsöl") will probably match poorly. To raise match rates, a Coolinary recipe page could offer an English version of the ingredient list with gram weights.
- Publishing every Coolinary recipe as a public page with schema.org markup is cheap to build and also helps SEO. For logging, though, it mostly reaches MFP-web (English), Cronometer Gold and MacroFactor users. It does not help Yazio, Lifesum or Fitatu users, who are most of the PL/DE target.

### Gaps
- Whether the MFP Recipe Importer is free or Premium in 2026, and whether it works in the MFP mobile app in DE/PL locales.
- No official statement from any tracker on JSON-LD or `nutrition` parsing, or on metric and non-English handling. Testing with real Polish and German pages is needed.
- Lose It! URL import status (official help page was 403).

## 3. Sharing and deep links (prefilled "add food")

### Takeaway
I found **no documented URL scheme, intent or deep link** in any of these trackers that opens "add food" with prefilled data from a third party. I also found no documented feature for sharing recipes or meals between users via link. For a tiny start-up this route is effectively unavailable unless a tracker is approached directly.

### Cited Findings
- The Lifesum DE App Store page shows only the store's generic "Teilen" button. It describes no in-app sharing of recipes. — [App Store DE](https://apps.apple.com/de/app/id286906691)
- The 2013 HelloFresh × Noom deal is a precedent for one-tap logging from food products: "each recipe kit delivered comes with a unique QR code that can be scanned so that key nutritional information is logged and updated in a user's personal Noom record" (2013, outdated). — [TechCrunch 2013](https://techcrunch.com/2013/04/03/hellonoom/amp/)
- FatSecret Brand Tools let a brand publish nutrition data into the FatSecret database, so users can find it through search or barcode. — [FatSecret Brands](https://platform.fatsecret.com/brands)

### Inferences
- The closest thing to "one-tap" logging available without a partnership is:
  1. Writing nutrition to HealthKit or Health Connect, for trackers that import it (MacroFactor daily totals; possibly Lose It!).
  2. FatSecret's diary API.
  3. Barcode or brand entries in tracker databases (FatSecret Brand Tools; MFP and Yazio user-contributed foods), which do not suit home-cooked dishes.
- A QR or deep-link handoff in the Noom style would require a bilateral deal.

### Gaps
- Undocumented app URL schemes (for example `mfp://`, `yazio://`, `lifesum://`) may exist but are not public or stable. I found no source.
- Whether MFP, Yazio or Lifesum allow sharing a user's custom recipe with another user via link in 2026.

## 4. Partner and brand programmes (brands, recipe partners, meal kits, insurers, employers)

### Takeaway
**MyFitnessPal** runs a paid ads and sponsored-recipe business (Sprouts, TABASCO, athlete collections). **Yazio's** public partner page offers only affiliate and content-creator programmes. **Lifesum** works through data partners (Consupedia), device makers (Samsung) and B2B benefit channels. **FatSecret** has Brand Tools and powers Samsung Health's barcode database. **Cronometer** partners with practitioner platforms. None of these is designed for a small recipe app to push its own recipes in for free.

### Cited Findings
- **MFP Sprouts (2021, outdated):** Sprouts became the first retailer partner and added 100+ recipes featuring Sprouts products. 2,169 Sprouts recipes had been used by app users. — [Sprouts press release](https://about.sprouts.com/press-release/sprouts-joins-myfitnesspal-as-first-retailer-partner/)
- **MFP × TABASCO (2024):** the MFP team reworked 10 recipes around TABASCO products. — [PR Newswire](https://prnewswire.com/news-releases/do-the-sauce-swap-myfitnesspal-and-tabasco-brand-team-up-for-spicy-summer-collab-302133113.html)
- **MFP × Cameron Brink (Jan 2026):** 8-recipe collection for Premium+ subscribers. — [press release mirror](https://business.mammothtimes.com/mammothtimes/article/gnwcq-2026-1-6-myfitnesspal-teams-up-with-basketball-all-star-cameron-brink-to-cut-through-new-year-nutrition-noise)
- **MFP ads business:** a pilot with 25 advertisers (Blue Apron, CVS, Target), positioned as health and wellness retail media. — [AdExchanger](https://www.adexchanger.com/marketers/myfitnesspal-wants-to-start-the-health-and-wellness-subsector-of-retail-media/); [Adweek](https://tunnel-sergio.adweek.com/commerce/myfitnesspal-is-launching-an-ad-business/)
- **Yazio:** the partner page lists only "Affiliate Partner" (website) and "Content Creator" (social). It claims "Over 95 million users" and "4.7 stars, over 1M ratings". — [yazio.com/en/become-a-partner](https://www.yazio.com/en/become-a-partner)
- Yazio's affiliate guideline (German) pays commission per YAZIO PRO purchase via a Yazio landing page. Contact: partnerships@yazio.com. — [Yazio terms PDF](https://filecontent.yazio.com/press/terms_and_conditions_de.pdf) (as summarised in search results)
- I found no Yazio partnership with a health insurer (Krankenkasse) or corporate-fitness provider. — (searches "Yazio Krankenkasse Kooperation", "Yazio Firmenfitness" returned nothing relevant)
- **Lifesum:**
  - Partnership with Consupedia (Apr 2025) on an AI-powered food database. — [Food and Drink Technology](https://www.foodanddrinktechnology.com/news/58547/lifesum-in-partnership-with-consupedia-to-build-the-worlds-most-informed-and-inclusive-ai-powered-food-database/); [Athletech News](https://athletechnews.com/consupedia-brings-new-layer-of-intelligence-to-lifesums-nutrition-app)
  - Amazon Halo nutrition coaching (2020, outdated) — [MobiHealthNews](https://www.mobihealthnews.com/news/swedens-lifesum-raises-67-million-food-tracking-app)
  - Samsung Galaxy Watch4 nutrition partner (2021, outdated) — [MobiHealthNews](https://www.mobihealthnews.com/news/samsung-inks-deal-nutrition-app-lifesum-new-galaxy-watch4)
- Lifesum appears to sell Premium through benefit platforms. Search results show partner landing pages at lifesum.com/gympass and lifesum.com/partner/edenred3, whose content was not verified. — [lifesum.com/gympass](https://lifesum.com/gympass); [lifesum.com/partner/edenred3](https://lifesum.com/partner/edenred3)
- **FatSecret × Samsung (Oct 2024):** Samsung Health's Food tracker added barcode scanning "powered by a partnership with fatsecret". — [Samsung Newsroom US](https://news.samsung.com/us/samsung-health-app-update-makes-accessing-health-records-managing-medications-food-tracking-easier)
- **Cronometer:** API partners are practitioner platforms (Everfit, Kalix, Practice Better). — [Cronometer Pro FAQ](https://cronometer.com/blog/cronometer-pro-faq)

### Inferences
- **MFP:** getting in means a paid media deal (retail media and sponsored recipes), which a bootstrapped start-up cannot afford. It is also US-centred.
- **Yazio:** the only self-serve programme is affiliate. Coolinary could earn commission by referring users to YAZIO PRO, but that gives no data integration.
- **Lifesum:** partnerships are mostly B2B, device or data deals. Cold outreach (contact@lifesum.com, per search results) is the only path, with low odds at pre-launch.
- **FatSecret Brand Tools:** a possible, unconventional route. Coolinary could publish its recipes as "brand" items so they are searchable in FatSecret and in Samsung Health's FatSecret-powered barcode/food search. Whether Brand Tools accepts non-packaged recipe items is unverified.

### Gaps
- Whether Yazio runs any brand, recipe, meal-kit or insurer partnerships in DACH. None were found in public sources.
- Formal criteria for MFP ads and recipe sponsorships, and EU availability.
- Brand Tools eligibility and pricing.

## 5. Apple Health / Health Connect: what each app writes and reads for nutrition

### Takeaway
Most trackers **write** nutrition to HealthKit or Health Connect (Yazio, Lifesum, MFP, FatSecret, Cronometer, partly Cal AI). Very few **read** nutrition from other apps. **MacroFactor** reads daily totals of calories, macros and micronutrients from Apple Health and Health Connect. **Lose It!** appears to import foods from Apple Health (help-page snippet only). **Yazio, Lifesum, Cronometer and FatSecret** read activity and weight but not nutrition. So Coolinary writing nutrition to the health platforms would reach only a minority of trackers.

### Cited Findings
- **Yazio (Apple Health):** "Apple Health transfers your activity data and body values, YAZIO transfers your activity data, body values and nutrients". From Apple Health it reads activities, workouts and weight, and uses Apple's active calories. — [Yazio help: YAZIO and Apple Health](https://help.yazio.com/hc/en-us/articles/360018833998-YAZIO-and-Apple-Health) (search snippet; page 403)
- Yazio has separate help articles for Health Connect, Samsung Health, Huawei Health and Fitbit sync. — [Yazio help: Health Connect](https://help.yazio.com/hc/en-us/articles/33129729688209-Syncing-Google-Fit-with-Health-Connect); [Yazio data sync overview](https://help.yazio.com/hc/en-us/articles/115003788369)
- **Lifesum (Apple Health):** exports food calories, food nutrition, water and weight. Imports exercise calories, steps and other activity. — [Lifesum help](https://help.lifesum.com/en/article/eagtau/); [App Store DE](https://apps.apple.com/de/app/id286906691)
- **Lifesum (Health Connect):** "Import: Exercises, Total calories burned, Steps tracked. Export: Nutrition (tracked food), Hydration (tracked water)". If nutrition is missing, the help page says to disable "Multimodal Tracking". — [Lifesum Helpshift FAQ 655](https://lifesum.helpshift.com/hc/faq/655-how-do-i-use-lifesum-with-health-connect---android)
- **MyFitnessPal:** the official Health Connect article covers nutrition syncing and manual sync (Menu > Sync). — [MFP support: Health Connect](https://support.myfitnesspal.com/hc/en-us/articles/10553948248973)
- Sahha (third party) says MFP writes calories and macros to Health Connect and Apple Health. — [Sahha](https://sahha.ai/integrations/myfitnesspal/) (secondary)
- Samsung Community users report MFP food intake does not sync to Samsung Health. — [Samsung Community](https://r1.community.samsung.com/t5/samsung-health/sync-mfp-food-intake-with-samsung-health/td-p/5295815)
- **Cronometer (Apple Health, post of 7 Mar 2024):** exports "13 Different Nutrients" plus dietary cholesterol. Imports activity, exercise, sleep, heart rate, weight and similar metrics, with **no nutrition import**. — [Cronometer blog: Apple Health](https://cronometer.com/blog/apple-health)
- A Cronometer staff member on the forum says it does not import nutrient information from Apple Health. — [Cronometer forum](https://forums.cronometer.com/discussion/6325/import-food-data-from-apple-health)
- Cronometer has a Health Connect help article (content not retrieved). — [Cronometer support](https://support.cronometer.com/hc/en-us/articles/22731903751316-Health-Connect)
- **MacroFactor (v2.9.3, 26 Aug 2024):** "MacroFactor can now import Calories, macros and micronutrients for the day from Apple Health and Health Connect". Weight is imported and exported. — [MacroFactor 2.9.3 release notes](https://macrofactor.com/version-2-9-3/)
- MacroFactor's source priority: "Manual, followed by Food Log, followed by Apple Health and Google Health Connect". A same-day entry in MacroFactor prevents the sync for that day. — [MacroFactor help: data sources](https://help.macrofactorapp.com/settings/data_sources/)
- **Lose It!:** syncing with Apple Health "is available for all iOS users". Foods written from Health into Lose It! get a heart icon in the log. There is no historical or backdated sync. — [Lose It help: Using Apple Health](https://loseit.zendesk.com/hc/en-us/articles/47773580085140-Using-Apple-Health-With-Lose-It) (search snippet; page 403)
- **FatSecret:** the official help page says the app "is fully integrated with Apple Health" and syncs calories burned and steps from Health into FatSecret. The Chinese-language version reportedly lists writing food/nutrition and weight to Health. — [FatSecret help](https://www.fatsecret.com/fatsecret-app-help/third-party-apps/apple-health-data-sync-with-fatsecret) (search snippet; page 403)
- **Samsung Health:** shares nutrition with Health Connect. — [Samsung dev blog](https://developer.samsung.com/health/blog/en/accessing-samsung-health-data-through-health-connect)
- **Cal AI:** the App Store listing says it syncs with Apple Health (HealthKit). — [App Store](https://apps.apple.com/us/app/-/id6755081229)

### Inferences
- **MacroFactor:** imports only **daily totals**, and only on days when the user has logged nothing in MacroFactor itself. It cannot mix Coolinary-delivered meals with manually logged snacks on the same day.
- **Lose It! (iOS):** may import individual food entries from Health. If confirmed, it is the best HealthKit target among the trackers covered here. Needs hands-on verification.
- **Yazio, Lifesum, Cronometer, FatSecret and MFP:** probably ignore nutrition that Coolinary writes to HealthKit or Health Connect. For MFP, whether it reads nutrition is unconfirmed (see Gaps).
- If Coolinary writes to HealthKit/Health Connect, it must not double count when the user also logs the same meal in a tracker that writes to Health. Health apps sum dietary energy across sources.

### Gaps
- Exact Yazio Health Connect nutrition behaviour, and whether Yazio reads any nutrition (help pages were 403).
- Whether MFP reads nutrition from Apple Health or Health Connect.
- Cronometer's exact list of the 13 exported nutrients, and whether dietary energy is included.
- Lose It! nutrition import details (official page blocked; snippet only).
- Whether Cal AI reads nutrition.

## 6. Unofficial APIs and Terms of Service

### Takeaway
Unofficial clients exist for MFP, Yazio, Cronometer, Lose It! and MacroFactor (many are 2025–26 MCP servers). They are fragile, and in MFP's case MFP has actively blocked them. MFP's ToS explicitly ban scraping and commercial use without a separate agreement. Building a product on these would be a legal and reliability risk, especially for an EU company handling health data under GDPR Art. 9.

### Cited Findings
- **MFP ToS (effective 18 Mar 2025):**
  - Bans "data mining or similar data gathering or extraction activities" "using automated means, including … robots, spiders, crawlers, scrapers".
  - Bans bypassing access-restriction measures.
  - Bans reverse engineering.
  - The Services are "intended only for your personal, non-commercial use", and users may not "collect, harvest, or use any Content, data … for commercial exploitation".
  — [MFP Terms of Service](https://www.myfitnesspal.com/terms-of-service)
- **python-myfitnesspal:** MFP added a hidden or invisible captcha to login, which stopped the library from logging in (issue #144). Later versions use `browser_cookie3` to borrow cookies from the user's browser session. Date not confirmed; likely before 2024. — [GitHub issue #144](https://github.com/coddingtonbear/python-myfitnesspal/issues/144); [Lightrun summary](https://lightrun.com/answers/coddingtonbear-python-myfitnesspal-login-now-includes-a-hidden-captcha-that-prevents-python-myfitnesspal-from-logging-in)
- An unofficial MFP MCP server README warns it "uses private endpoints intended for their own clients" and may violate the ToS. — [Glama: myfitnesspal-mcp](https://glama.ai/mcp/servers/thesaaspreneur/myfitnesspal-mcp/tree)
- **Yazio:** `yazio-go-cli` relies on YAZIO's "unofficial/private API", where endpoints may change without notice. The `yazio-mcp` server says Yazio provides no official API. — [GitHub itzptk/yazio-go-cli](https://github.com/itzptk/yazio-go-cli); [mcpservers.org yazio-mcp](https://mcpservers.org/servers/github-com-fliptheweb-yazio-mcp)
- **Cronometer:** `cronometer-mcp` uses the web app's GWT-RPC protocol. `cronometer-api-mcp` uses the mobile app's REST endpoints. `cronopy` also exists. — [PyPI cronometer-mcp](https://pypi.org/project/cronometer-mcp/0.1.4/); [Socket cronometer-api-mcp](https://socket.dev/pypi/package/cronometer-api-mcp); [PyPI cronopy](https://pypi.org/project/cronopy/)
- **Lose It!:** `loseit-mcp` logs food through private web endpoints and "can break whenever they ship a new web build". — [Glama loseit-mcp](https://glama.ai/mcp/servers/cabird/loseit-mcp/tree)
- **MacroFactor:** the `sjawhar/macrofactor` client and MCP server are archived and "no longer functional". Firebase App Check tokens "can only be minted by the official iOS app (App Attest), Android app (Play Integrity), or a Firebase-registered web client", which blocks third-party clients. — [GitHub sjawhar/macrofactor](https://github.com/sjawhar/macrofactor)

### Inferences
- Unofficial routes would mean storing user credentials or session cookies for a third-party health service. That violates MFP's ToS outright and probably others'. It also creates GDPR special-category-data exposure. It is not viable for a commercial product, and at most suitable for a founder's personal prototype.
- The trend is toward stronger attestation (App Check, hidden captcha), so these routes will keep getting harder.

### Gaps
- I did not review the ToS of Yazio, Lifesum, Cronometer, Lose It! or MacroFactor for specific anti-scraping clauses.
- No documented enforcement actions (legal letters, account bans) against third-party developers, beyond MFP's captcha change.

## 7. Market presence in Poland and Germany, and Premium prices

### Takeaway
In Poland, **Fitatu** dominates, with **FatSecret** and **MyFitnessPal** usually named next and **Yazio** visible in the charts (#30 in Health & Fitness on the PL App Store). In Germany, MFP and Lifesum each have over 80k App Store ratings but rank around #74–84 in Health & Fitness. Yazio claims more than 95M users globally. Premium costs roughly €50–100 per year, or 130–390 PLN per year.

### Cited Findings
App Store snapshots, 30 Oct 2026:

| App | Store | H&F rank | Ratings | Listed prices |
|---|---|---|---|---|
| Yazio ("AI Calorie Counter by Yazio") | PL | #30 | 33k (4.7) | 44.99–139.99 zł (e.g. "12 Months" 64.99 or 129.99 zł; "Semi-Annual" 49.99 or 84.99 zł; "3 Months" 64.99 zł). Polish UI. — [App Store PL](https://apps.apple.com/pl/app/id946099227) |
| MyFitnessPal | PL | #107 | 11k (4.7) | Monthly 34.99 or 97.99 zł; yearly 184.99 or 389.99 zł. Polish UI. — [App Store PL](https://apps.apple.com/pl/app/id341232718) |
| MyFitnessPal | DE | #84 | 80,310 (4.5) | Monthly 9.99 or 21.99 €; yearly 49.99 or 87.99 €. German and Polish UI. — [App Store DE](https://apps.apple.com/de/app/id341232718) |
| Lifesum | DE | #74 | 81,900 (4.5) | 1 Monat 9.99–14.99 €; 3 Monate 21.99–39.99 €; 1 Jahr 49.99–99.99 €. Polish added recently. — [App Store DE](https://apps.apple.com/de/app/id286906691) |

- **Yazio globally:** "Over 95 million users", 4.7 stars from more than 1M ratings. — [Yazio partner page](https://www.yazio.com/en/become-a-partner)
- **Fitatu:** about 1.2M monthly active users (Fitatu data cited by Telepolis). — [Telepolis](https://www.telepolis.pl/tech/aplikacje/fitatu-polska-aplikacja-odchudza-zdjecie-posilku-ozempic)
- Polish online rankings mainly cite Fitatu, FatSecret and MyFitnessPal. — [Geekweek/Interia](https://geekweek.interia.pl/lifestyle/news-jaka-aplikacja-do-liczenia-kalorii-jest-najlepsza,nId,6402954)
- Komputronik's 2026 guide ranks Fitatu #1 and FatSecret #2 among free apps. — [Komputronik](https://nano.komputronik.pl/n/aplikacje-do-liczenia-kalorii-najlepsze-kalkulatory-kalorii-na-smartfon/)
- **Cal AI:** more than 15M downloads globally. — [The Next Web](https://thenextweb.com/news/myfitnesspal-acquires-cal-ai-the-viral-calorie-tracking-app-built-by-teens)

### Inferences
- Several price points per product suggest tiered plans (MFP Premium vs Premium+) and introductory offers. The higher MFP figures (21.99 €/month, 87.99 €/year; 389.99 zł/year) probably correspond to Premium+, which is not labelled on the store page.
- For Coolinary's PL launch, the trackers that matter after Fitatu are FatSecret, MFP and Yazio. In DE, Yazio is likely the most relevant, given it is a German company and claims over 95M users. I did not get a DE rank for Yazio.
- Cronometer, Lose It! and MacroFactor are niche or US-centric in PL/DE (inference; no PL/DE rank data found). Integrating with them gives little reach.

### Gaps
- Yazio's DE App Store rank and DE price, plus PL/DE ranks for FatSecret, Cronometer, Lose It!, MacroFactor and Cal AI.
- Google Play download counts per country, which matter because Android has a larger share in Poland.
- Third-party download and MAU data (Sensor Tower, data.ai) for PL/DE (not found in free sources).
- Lose It! and Cronometer EUR prices.

## 8. Per-app practicality matrix for a tiny EU start-up (synthesis)

### Takeaway
The only official route that both reads and pushes for a bootstrapped company is **FatSecret** (free Premier tier; diary writes via 3-legged OAuth; PL/DE data subject to approval). Everything else is HealthKit or Health Connect writes (reaching only MacroFactor and possibly Lose It!), schema.org recipe pages (reaching MFP-web-English, Cronometer Gold and MacroFactor), or business development with low odds of success.

### Cited Findings
- All cells in the matrix below are backed by the citations in sections 1–7. No new facts are introduced here.

### Inferences

| App | Official API (2026) | Can Coolinary PUSH meals? | Can Coolinary READ intake? | URL recipe import | Reads nutrition from HealthKit / Health Connect? | Practicality for Coolinary |
|---|---|---|---|---|---|---|
| FatSecret | Yes, self-serve: Basic / Premier Free / Premier | **Yes**: `food_entry.create` (OAuth 1.0 3-legged profile); needs FatSecret food IDs | Yes: profile diary methods | Not found | No: reads activity and steps only (snippet) | **High.** Free for <$1M revenue start-ups; ask for PL/DE localisation; test consumer-app diary visibility |
| MyFitnessPal | Closed ("not accepting requests") | No (only via a deal) | No | Yes; English website only; Premium status disputed | Unconfirmed | **Low.** Schema.org pages for English users only; paid ads route |
| Yazio | None (public API request "under consideration") | No | No | No (feature request) | No (reads activity and weight; writes nutrients) | **Low technically.** Affiliate programme only; worth BD outreach given DE strength |
| Lifesum | None public | No | No | Not found | No (reads exercise and steps; writes nutrition) | **Low.** B2B/BD only |
| Cronometer | Select partners only (practitioner platforms) | No | Partner-only | Yes, Gold-only; USDA/NCCDB matching | No | **Low** (niche in PL/DE) |
| Lose It! | Beta partner API, read-only, business only | No (no write access) | **Yes**, if accepted (OAuth 2.0 + webhooks) | Unclear | Possibly yes on iOS (snippet: foods from Health shown with a heart icon) | **Medium-low.** Small EU user base (inferred) |
| MacroFactor | None (2023 statement; App Check blocks clients) | Indirect: write daily totals via HealthKit/HC | No | Yes ("Import from Link") | **Yes**: daily calories, macros and micros; only on days without a MacroFactor food log | **Low-medium.** Niche; all-or-nothing per day |
| Samsung Health | Official Data SDK; Nutrition write needs an approved partner request | **Yes** if approved (or via Health Connect sync) | Yes if registered | n/a | n/a (it is the store) | **Medium.** Android-heavy PL market; approval criteria unknown |
| Cal AI (MFP-owned since Mar 2026) | None found | No | No | n/a | Unknown (writes to HealthKit per listing) | **Low** |

- Recommended low-cost stack (inference):
  1. Write per-meal nutrition to HealthKit and Health Connect (covered by the other researcher).
  2. Publish each recipe as a public page with schema.org/Recipe markup, including gram weights and an optional English ingredient list.
  3. Prototype a FatSecret Premier Free integration for one-tap diary logging and an EU-localised food database.
  4. Treat Yazio, Lifesum and MFP as business-development targets after there is traction, not as integrations for launch.

### Gaps
- No hands-on test was done of any importer with Polish or German schema.org pages.
- No hands-on test of whether FatSecret API diary entries appear in the consumer fatsecret app.
- Samsung partner-request acceptance criteria are not known.
