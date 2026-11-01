# Fitatu: integration surface, partner programmes, legal/commercial realism, company facts (as of 30 Oct 2026)

Scope note: research done on 30 Oct 2026. Items older than 2024 are flagged **[pre-2024, possibly outdated]**. "Help-centre" pages on fitatu.com carry no visible date. Where I write "inference", it is my reasoning and has no direct source.

## 1. Does Fitatu offer any public or partner API, B2B offer, dietitian panel or white-label?

### Takeaway
No. I found no public API, developer portal, partner API, white-label product or dietitian/B2B panel. Fitatu's published partnership channel is marketing-only: an e-mail address for brand, advertising, influencer, food-producer and event-sponsorship requests. Every programmatic access route I found is unofficial and reverse-engineered.

### Cited Findings
- Fitatu's help page "Jakie są możliwości współpracy z Fitatu®?" lists only these partnership types: influencer, food producer adding products to the database, buying an in-app ad campaign, marketing partnership, and event sponsorship, plus a jobs link. Contact is **wspolpraca@fitatu.com**; customer service is help@fitatu.com. It does not mention an API, integrations, technology partners, caterers, restaurants, dietitians or B2B. — [Fitatu help: możliwości współpracy](https://www.fitatu.com/pl/help/jakie-sa-mozliwosci-wspolpracy-z-fitatu)
- The English help centre has a matching entry: "I'd like to start a marketing collaboration with Fitatu®". — [Fitatu help EN](https://www.fitatu.com/en/help/i'd-like-to-start-a-marketing-collaboration-with-fitatu.)
- The Terms (effective 13.07.2026) never use the words "API" or "partner" ("partner" was searched for in the terms text). — [Fitatu Regulamin](https://www.fitatu.com/pl/regulations?p=959)
- The only Python client I found calls itself "Unofficial Python API client for Fitatu" and says it is "not affiliated with or endorsed by Fitatu Sp. z o.o.". It was "reverse-engineered" from the web and Android apps. — [GitHub Jezue/fitatu_library](https://github.com/Jezue/fitatu_library)
- Dietitians: a 2023 trade-show page pitched Fitatu as a tool for trainers and dietitians. The mechanism was clients exporting their data with Premium, not a dedicated panel. **[pre-2024, possibly outdated]** — [FIWE: The Fitatu® app as a tool for trainers and dietitians](https://fiwe.pl/en/news/the-fitatu-app-as-a-tool-for-trainers-and-dietitians/)
- Data export (XLS/CSV) is a Premium feature in the plan table. — [Fitatu Regulamin §III.1](https://www.fitatu.com/pl/regulations?p=959); [Google Play listing](https://play.google.com/store/apps/details?id=com.fitatu.tracker&hl=en_US)
- The Terms say Fitatu's content "is not dietetic services" and advise users to consult a doctor or dietitian (§IX.4(b)). — [Fitatu Regulamin](https://www.fitatu.com/pl/regulations?p=959)
- Ownership coverage (Feb 2026) discusses AI, international expansion, possible M&A and new healthy-living categories "not before about a year from now". It does not mention B2B, partnerships or integrations. — [XYZ.pl, 3 Feb 2026](https://xyz.pl/tworcy-sympatii-oddali-kontrole-w-fitatu-grupa-znanych-inwestorow-rozkreci-zywieniowa-aplikacje-w-europie-zachodniej/)

### Inferences
- Coolinary has no self-serve route. Any official integration would have to be negotiated through wspolpraca@fitatu.com, and that inbox is set up for marketing deals, not technical ones.
- The new investor-controlled management (CEO since early 2026) is focused on DACH/Western-Europe expansion and subscriber growth. A tiny pre-launch partner asking for API access is unlikely to be a priority. A paid content or co-marketing deal is more plausible than an API.

### Gaps
- No public statement says whether Fitatu would ever offer a partner API. Only a direct enquiry would answer this.
- I could not check for an unlisted dietitian or "Fitatu for Business" product. Search returned none.

## 2. Existing partner integrations and how a company becomes a partner (catering, restaurants, recipes, retailers, brands, sponsored content)

### Takeaway
Fitatu partnerships are content and advertising deals: products added to the food database, sponsored recipes in the recipe catalogue, banners, and co-branded Premium discounts (Lidl Plus, InPost, MediaMarkt). Restaurant and retailer items appear as **database entries**, not live menu feeds. I found **no official diet-catering menu import**. The only catering-to-Fitatu bridges are unofficial hobby tools that use the private API. Prices are not public.

### Cited Findings
- **Food producers:** they e-mail wspolpraca@fitatu.com stating how many products to add. The page does not say whether this is paid or what data format is needed. — [Fitatu help: Jestem producentem żywności](https://www.fitatu.com/pl/help/jestem-producentem-zywnosci-i-chce-dodac-moje-produkty-do-bazy-fitatu.)
- **Ad and partner formats** sold through agency Cityboard Media (article dated 20 Aug 2024):
  - Branded promotion: a partner logo on a Premium discount campaign, about 500k in-app impressions and 200k pushes over 2 days.
  - Planner text banner for free users: 23M+ impressions and 234k+ unique users in 30 days.
  - In-app banners (468×60 top, 320×480 full screen) targeted by device, gender, age and city.
  - Social activations.
  - Product-database work: adding products, cleaning data, completing Fitatu Score.
  - **"Recipe database integration": recipes built on the advertiser's products are placed permanently in the app.**
  - Pricing "on request" via Cityboard (office@cityboard.pl). The article states about 1M monthly users in Poland.
  - — [Cityboard: Reklama na Fitatu](https://cityboard.pl/blog/713-reklama-na-fitatu)
- **Retail co-marketing (Jan 2026):** Lidl Plus users got a 30% code for an annual Premium+AI plan via the Lidl Plus "Benefit Plus" section. It was redeemed on fitatu.com and valid until the end of January. There was no recipe or product integration. — [Fitatu Blog, 15 Jan 2026](https://www.fitatu.com/blog/nowy-rok-nowe-cele-zrob-pierwszy-krok-do-formy-z-fitatu-i-lidl-plus/)
- Other co-promotions appear only as deal-site titles, which I did not open:
  - "Premium Fitatu 35 zł taniej wymieniając InCoiny w InPost" (InPost loyalty points, summer 2026). — [Pepper.pl](https://www.pepper.pl/promocje/premium-fitatu-35zl-taniej-wymieniajac-incoiny-w-inpost-1256647)
  - A myMediaMarkt-only Fitatu offer. — [Pepper.pl](https://www.pepper.pl/promocje/fitatu-w-cenie-na-ktora-czekalem-1155355)
- **Restaurants and retailers as database content:**
  - Polish site: dishes from McDonald's, KFC, Burger King, Subway, Pizza Hut and Starbucks; store-brand products from Biedronka, Lidl, Kaufland. — [fitatu.com/pl](https://www.fitatu.com/pl)
  - English Play listing: Tesco, Asda, Morrisons, Sainsbury and Lidl products. — [Google Play](https://play.google.com/store/apps/details?id=com.fitatu.tracker&hl=en_US)
- **Recipe catalogue:** two categories, recipes made by Fitatu and community recipes added by users and "verified by our dietitian team" before publication. — [Fitatu help: Co znajdę w "Przepisach"?](https://www.fitatu.com/pl/help/co-znajde-w-%22przepisach%22)
- **"Wybór menu" (Premium):** 7 ready-made menus (Balans, Bezglutenowy, Bezlaktozowy, Keto, Mniej cukru, Wege, Wysokobiałkowy) built from catalogue recipes.
  - Each suggestion has a calorie range with a ±100 kcal buffer, and users fit meals to their daily goal themselves.
  - The page names no partners or caterers and says nothing about external menu import. Users can add their own dishes to the plan.
  - — [Fitatu help: "Wybór menu"](https://www.fitatu.com/pl/help/%22wybor-menu%22-jak-dziala-ta-funkcja)
- **Catering:** searches for an official Fitatu integration with Maczfit, Kuchnia Vikinga or Dietly found none. What exists:
  - **VikFit**: a free desktop app (© 2025, by "Wavelski") that pulls Kuchnia Vikinga orders' calories and macros into the user's Fitatu diary.
    - It says it is "narzędziem nieoficjalnym" (an unofficial tool), may break if either service changes its API, and stores passwords locally in Keychain.
    - It was inspired by an open-source script by "mzielinski".
    - — [vikfit.pl](https://vikfit.pl/); [VikFit FAQ](https://vikfit.pl/faq/)
  - The unofficial `fitatu-api` library ships a `catering_sync.py` example (PR #21). It imports a caterer's menu (JSON/CSV: date, meal, name, kcal/macros or grams + barcode/search phrase) into the Fitatu diary without duplicates. — [GitHub PR #21](https://github.com/Jezue/fitatu_library/pull/21)
  - `meal-optimizer`, an open-source tool, picks meals across Maczfit, Kuchnia Vikinga and Zdrowy Catering (Dietly platform). It does not connect to Fitatu. — [GitHub meal-optimizer](https://github.com/MarekSwiechowicz/meal-optimizer)
  - A Dietly app user asks for Apple Health integration because "Cateringi wrzucają posiłki i makroskładniki" (caterers post meals and macros). This suggests caterers' macro data lives in Dietly, not in Fitatu. — [Dietly App Store (PL)](https://apps.apple.com/pl/app/dietly/id1642713520)
- Fitatu's channel reach as of Aug 2024: Facebook page 104k, Facebook group 107k, TikTok 110k, Instagram 15k. — [Cityboard](https://cityboard.pl/blog/713-reklama-na-fitatu)

### Inferences
- The closest official equivalent to "push Coolinary recipes into Fitatu" is the paid **sponsored "recipe database integration"**: Fitatu places partner recipes permanently in its catalogue. It is designed for food brands that pay for placement, not for a startup syndicating thousands of dynamic, per-person-scaled recipes.
- Polish users who want catering meals in Fitatu today resort to unofficial scripts. That signals unmet demand Fitatu has not addressed officially, but also that Fitatu does not run a catering-import programme Coolinary could join.

### Gaps
- No public price list exists for any partner or ad format. Cityboard's 2024 figures are agency marketing claims.
- I could not confirm whether any caterer has an official Fitatu deal. Caterers' sites and Fitatu pages showed none.
- I could not confirm whether adding products to the database is free or paid.

## 3. Health-platform integrations (Apple Health/HealthKit, Health Connect, Google Fit, Samsung, Garmin, Fitbit, etc.): what does Fitatu read and write?

### Takeaway
Fitatu documents two-way sync with Apple Health and Health Connect:
- It **writes meals, weight and height**. Writing data to Apple Health and Health Connect is a **Premium-only** feature.
- It **reads only burned calories, workouts and activity**, plus weight from scales.

Nothing documents Fitatu **reading nutrition written by other apps**. So meals Coolinary writes to HealthKit or Health Connect would very likely not appear in a Fitatu diary. The reverse is a clean, ToS-safe READ route: for **Premium** users, Coolinary could read Fitatu's meal data from HealthKit or Health Connect.

### Cited Findings
- Fitatu help (EN, undated) on Apple Health and Health Connect:
  - "two-way sync"; it receives "data on burned calories, workouts, and activity" and sends "meals, weight, and height".
  - It receives directly from Huawei Health, Samsung Health (Android), Strava, Garmin Connect and Polar Flow.
  - "To avoid duplicate data, it's best to use only one sync path."
  - — [Fitatu help: Can Fitatu® sync with other devices and apps?](https://www.fitatu.com/en/help/can-fitatu-sync-with-other-devices-and-apps)
- The Polish version says the same: "odbieranie danych o spalonych kaloriach, treningach i aktywności" (receiving burned calories, workouts and activity) and "wysyłanie danych o posiłkach, wadze i wzroście" (sending meals, weight and height). Consumed calories, macros, steps and reading other apps' nutrition are not mentioned. It also points users to in-app "Ustawienia > Powiązane aplikacje" (Settings > Connected apps) for per-integration details. — [Fitatu help PL](https://www.fitatu.com/pl/help/czy-jest-mozliwa-synchronizacja-fitatu-z-innymi-urzadzeniami-i-aplikacjami); [Fitatu help: jak działają poszczególne synchronizacje](https://www.fitatu.com/pl/help/jak-dzialaja-poszczegolne-synchronizacje)
- Plan table in the Terms (§III.1), which is authoritative on tiers:

  | Feature | Free | Premium |
  |---|---|---|
  | "Synchronizacja z Health Connect" (import) | Yes, app only | Yes |
  | "Wysyłka danych do Health Connect i Apple Zdrowie" (sending data) | No | Yes, app only |
  | "Synchronizacje z pozostałymi aplikacjami fitness" | No | Yes |
  | "Eksport danych" | No | Yes |

  — [Fitatu Regulamin §III.1](https://www.fitatu.com/pl/regulations?p=959)
- The Fitatu homepage advertises two-way sync with Health Connect and Apple Health ("100+" devices and apps), plus Huawei Health, Samsung Health, Strava, Garmin Connect, Fitbit and Polar Flow as "receive data only". — [fitatu.com/pl](https://www.fitatu.com/pl)
- The Google Play listing (updated 14 Oct 2026) says:
  - "Integration with Google Fit, Garmin Connect, Fitbit, Samsung Health, Huawei Health, and Strava".
  - "Import activity data from Adidas Running (Runtastic) and Zepp Life via Google Fit".
  - "Additional backup & export options, including **sending nutrition and weight data to Google Fit**".
  - — [Google Play](https://play.google.com/store/apps/details?id=com.fitatu.tracker&hl=en_US)
- Store descriptions differ from the help centre. The PL and DE App Store descriptions name Google Fit, Garmin Connect, FitBit, Samsung Health, Huawei Health and Strava, and do not name Apple Health. The help centre and Terms do name Apple Health ("Apple Zdrowie"). — [App Store PL](https://apps.apple.com/pl/app/fitatu-licznik-kalorii-i-dieta/id1011095795?l=pl); [App Store DE](https://apps.apple.com/de/app/id1011095795)
- An Apple Watch app exists (watchOS 10+). It shows daily calorie goal and macros. — [App Store US](https://apps.apple.com/us/app/id1011095795)
- A separate help article covers sending body weight from smart scales or synced apps. — [Fitatu help: masa ciała z inteligentnych wag](https://www.fitatu.com/pl/help/czy-moge-przesylac-informacje-o-masie-ciala-z-inteligentnych-wag-lub-zsynchronizowanych-aplikacji-do-waszej-aplikacji)
- The Samsung Community thread contradicts this. A user says Fitatu support told them a Samsung Health integration was not planned, while other users report sync via Health Connect. The thread is undated in my capture and likely older; current Fitatu pages list Samsung Health. — [Samsung Community](https://eu.community.samsung.com/t5/aplikacje-i-us%C5%82ugi/samsung-health-synchronizacja-z-fitatu/td-p/1289354)
- A sponsored Tabletowo article (Mar 2025) says Fitatu syncs activity data with Google Fit and Apple Health to account for calories burned. — [Tabletowo, 21 Mar 2025](https://www.tabletowo.pl/fitatu-premium-ai-szacowanie-kalorii-ze-zdjec/)

### Inferences
- **WRITE route (Coolinary → Health → Fitatu): effectively closed.** Fitatu documents reading only activity, burned calories and weight. A Coolinary meal written to HealthKit or Health Connect as dietary energy or macros would almost certainly not be imported into the Fitatu diary. It could also be double-counted by other apps.
- **READ route (Fitatu → Health → Coolinary): open but partial.**
  - Only Fitatu **Premium** users export meals to HealthKit or Health Connect, roughly 240k paying of about 1.2M MAU (§6).
  - Coolinary could request read permission for dietary energy and macro types and show "already eaten today" for those users. This is fully within Apple and Google's frameworks and does not touch Fitatu's ToS.
  - Weight is written too.
  - Calorie **goals** and targets are not documented as exported, so Coolinary would still need its own target or the user's manual input.
- Meals exported to HealthKit are typically written as aggregate nutrition samples, not dish names or recipes. Coolinary would get totals, not item-level diary detail. This is unverified for Fitatu specifically.

### Gaps
- I found no source listing the exact HealthKit or Health Connect data types Fitatu requests: dietary energy only versus full macros/micros, per-meal or per-day granularity, and whether goals are written. This needs a hands-on test on iOS (Health app > Sources > Fitatu) and Android (Health Connect > App permissions > Fitatu).
- No explicit statement either way on whether Fitatu reads nutrition samples written by other apps. The inference that it does not rests on the documented scope (activity, burned calories, weight).
- I did not find an official Google Health Connect "compatible apps" list naming Fitatu.

## 4. Recipe creation, import and sharing in Fitatu, and deep links

### Takeaway
Users create their own dishes and recipes manually inside Fitatu. There is **no import from a URL, QR code or another app**. Fitatu has a built-in **"Udostępnij" (share) for meals** that produces a link the recipient can add to their own day and edit. Fitatu's iOS universal-link and Android app-link files show the app opens `fitatu.com` share, recipe, product, search and AI-estimation paths, so a third party can deep-link users to existing Fitatu content but cannot create content through a link.

### Cited Findings
- Recipes: users can use existing suggestions or "add their own recipe, which is saved to their account". Recipes can be favourited, "shared with family or friends", rated or commented on. Community recipes are verified by dietitians before being shared publicly. — [Fitatu help: Co znajdę w "Przepisach"?](https://www.fitatu.com/pl/help/co-znajde-w-%22przepisach%22)
- Meal sharing:
  - Tap the three dots next to a meal, then "Udostępnij". Meals go to family and friends via messengers or social media, and the link can be saved in e.g. Notes.
  - "The person who receives the meal can add it to a day and meal of their choice and edit it."
  - The page does not say whether the recipient needs an account or Premium.
  - — [Fitatu help: Czy mogę udostępniać posiłki?](https://www.fitatu.com/pl/help/czy-moge-udostepniac-posilki)
- Older guides describe own dishes under the "WŁASNE" tab of the search screen and a "Udostępnij potrawę" (share dish) step when creating a dish. Wykop users complained that shared recipes were hard for a partner to find, and that public visibility needs submission and approval. **[pre-2024, possibly outdated]** — [joannahasnik.wordpress.com (2019)](https://joannahasnik.wordpress.com/2019/11/09/jak-dodac-wlasna-potrawe-w-fitatu/); [Wykop thread](https://wykop.pl/wpis/39090839/jak-to-jest-w-fitatu-z-udostepnianiem-swoich-posil); [annapacholak.pl](https://annapacholak.pl/fitatu-jak-dodac-wlasna-potrawe-i-oszczedzic-czas-czesc-3/)
- **Universal links (iOS)**, from `https://www.fitatu.com/.well-known/apple-app-site-association` (fetched 30 Oct 2026). App ID `4TSPU58Q6A.com.fitatu.tracker` handles these paths:
  - `/app/share/*` and `/??/share/meal/*`, where `??` looks like a two-letter language code.
  - `/app/recipe/*`, `/app/product/*`, `/recipe`, `/product`.
  - `/app/search`, `/app/search/barcodeScanner`, `/app/search/ai`.
  - `/app/menu`, `/app/menu-plans`, `/app/menu-plan/*`, `/app/all-recipes`.
  - `/app/photo-ai-estimation`, `/app/voice-ai-estimation`, `/app/text-ai-estimation`.
  - `/app/fitness-app-auth/*`, `/app/premium-offer`, `/app/special-promotion-access-group/*`, among others.
  - — [Fitatu AASA file](https://www.fitatu.com/.well-known/apple-app-site-association)
- **Android App Links:** `https://www.fitatu.com/.well-known/assetlinks.json` grants `handle_all_urls` (and `get_login_creds`) to package `com.fitatu.tracker` and the `com.fitatu.tracker.uat` test build. — [Fitatu assetlinks.json](https://www.fitatu.com/.well-known/assetlinks.json)
- Search found no indexed public examples of `fitatu.com/.../share/meal/...` links. — (search, 30 Oct 2026; no result)
- **No URL or recipe import:** a search for a Fitatu recipe-from-link or AI recipe-import feature found nothing. Competitors such as MacroFactor have "Import from Link". — [MacroFactor help](https://help.macrofactorapp.com/en/articles/259-import-recipes-from-link)
- Fitatu does have AI **text**, **voice** and **photo** meal estimation, which parse free text, speech or a photo into database products with portions (Premium+AI, up to 600 uses per month each) — [Fitatu Regulamin §III.1](https://www.fitatu.com/pl/regulations?p=959); [Google Play](https://play.google.com/store/apps/details?id=com.fitatu.tracker&hl=en_US)
- An open-source library's docs say `get_recipe` returns **another user's private recipe by numeric id**, including the author's username and userId. This is reported by the library's own API reference, not independently verified. — [fitatu_library API_REFERENCE.md](https://github.com/Jezue/fitatu_library/blob/main/docs/API_REFERENCE.md)

### Inferences
- "One-tap push of a Coolinary recipe into Fitatu" is **not achievable officially**. Links can only open content that already exists in Fitatu: a catalogue recipe ID, a product, or a meal someone shared.
- Practical near-one-tap options:
  - **(a)** If Coolinary recipes were placed in the Fitatu catalogue via a paid content deal, Coolinary could deep-link to `/app/recipe/{id}`, and the user adds it with their own portion size inside Fitatu.
  - **(b)** For Premium+AI users, Coolinary could offer a "Copy for Fitatu" button that copies a plain-text ingredient list with grams ("150 g chicken breast, 80 g rice…"). The user opens Fitatu's text-AI estimation (`/app/text-ai-estimation` is a universal-link path) and pastes. Fitatu's AI matches the items to its database. This needs no API or ToS breach, but works only for about 150k Premium+AI users and needs a paste step.
  - **(c)** A Coolinary-owned Fitatu account could create meals and publish share links. These have fixed portions (no per-person scaling), an unknown lifetime, and are ToS-grey for commercial distribution (see §5).
- Whether the deep links accept query parameters, for example to pre-fill text-AI search, is unknown. It would need testing.

### Gaps
- No official documentation of the share-link URL structure, link lifetime, whether non-users can open it, or whether recipes (not only meals) get public links.
- Not tested whether `/app/text-ai-estimation` or `/app/search` accept a query string to pre-fill input.

## 5. Unofficial/reverse-engineered APIs; what the Terms say; legal risk for a commercial app

### Takeaway
One active open-source client exists: `fitatu-api` by Jezue (alpha, 21 stars, v0.6.0 in Oct 2026). It can read and write almost everything, including the diary, private recipes, products, goals and weight. Technically it would deliver both READ and PUSH. But Fitatu's Terms (effective 13.07.2026) **expressly ban automated access with bots and scripts, circumventing technical limits, decompiling, collecting or processing data for other purposes, and letting third parties use an account**. A commercial app that holds users' Fitatu credentials and calls the private API would breach the ToS, process GDPR special-category data, and depend on an API that can change without notice. For a bootstrapped startup this is high-risk and not advisable.

### Cited Findings
**Unofficial clients**
- `Jezue/fitatu_library` (PyPI `fitatu-api`), MIT, Python 3.11+, status badge "alpha". At the time fetched: 21 stars, 1 fork, 91 commits. — [GitHub](https://github.com/Jezue/fitatu_library)
  - Capabilities: diary read and write (add products, recipes, custom rows; move, copy, delete); activities, weight, measurements and water; goals and settings (read and set targets, macro split); summaries and export; food search; private `create_product`, `create_recipe` and `create_food`; 47 LLM tools and an MCP server.
  - Shared-catalogue writes (`create_recipe(shared=True)`, product proposals) are excluded "on purpose".
  - Auth: e-mail and password, or bearer and refresh tokens copied from the web session for Apple/Facebook-login accounts. The bearer token lasts about 1 hour.
  - Hosts: `https://pl-pl.fitatu.com/api`, `https://de-de.fitatu.com/api`.
  - It sends a fixed web-app client profile (app version, API key and secret). It ships a `doctor` command, a contract-check script and a weekly live test "to detect API drift".
  - Route catalogue taken from Android app 4.17.1.
  - Some endpoints return 403 without Premium.
  - Disclaimer: automated access "may break Fitatu's Terms of Service … use at your own risk".
- v0.6.0 was released to PyPI in mid-October 2026. The project is listed in MCP and agent directories. — [PyPI fitatu-api 0.6.0](https://pypi.org/project/fitatu-api/0.6.0/); [GitHub release v0.6.0](https://github.com/Jezue/fitatu_library/releases/tag/v0.6.0); [mcpservers.org](https://mcpservers.org/servers/jezue/fitatu_library)
- VikFit (Kuchnia Vikinga → Fitatu) is unofficial and warns it may stop working if either API changes. — [vikfit.pl](https://vikfit.pl/)
- `jmolinski/respo2fitatu` appeared in search results but returned 404 on 30 Oct 2026 (removed or private). — [GitHub (404)](https://github.com/jmolinski/respo2fitatu)
- No Home Assistant integration for Fitatu was found. — (search, 30 Oct 2026)

**Fitatu Terms (Regulamin), effective 13.07.2026 (§XVII.4)** — [Fitatu Regulamin](https://www.fitatu.com/pl/regulations?p=959)
- §VI.4(a) bans using the Services "w sposób zautomatyzowany, z wykorzystaniem botów, skryptów lub innych narzędzi automatyzujących" (in an automated way, using bots, scripts or other automation tools).
- §VI.4(b) bans circumventing "ograniczeń technicznych lub funkcjonalnych" (technical or functional limitations).
- §VI.4(c) bans use of one account by more than one user.
- §VI.4(d) bans disrupting the Services.
- §V.5 bans "dekompilacji i dezasemblacji" (decompiling and disassembling) and "odtwarzania jego kodu źródłowego" (reconstructing its source code).
- §IX.3 bans using the services and published data for purposes other than using Fitatu's services, including "kopiowanie, gromadzenie, przetwarzanie i upowszechnianie danych" (copying, collecting, processing and distributing data).
- §IV.5–6: users must not disclose their password. "Korzystanie z Kont innych Użytkowników lub udostępnianie Konta osobom trzecim jest zabronione." (Using other users' accounts or sharing an account with third parties is prohibited.)
- §V.1: all rights to the Platforms, "w całości i we fragmentach" (in whole and in part), including copyright and related rights, belong to Fitatu. §V.4: passing content to third parties is allowed only "przy wykorzystaniu narzędzi zawartych w Platformach" (using the tools built into the Platforms). §V.6 bans copying, selling or otherwise putting the app or content into circulation.
- §XIV.9(c)–(d) bans user materials that serve "działalności konkurencyjnej wobec Fitatu" (activity competitive with Fitatu) or unauthorised advertising.
- **§XIV.15–17: users grant Fitatu a free, transferable, non-exclusive licence, unlimited in time and territory, to "Materiały", explicitly including recipes and photos, for purposes including "trenowania i ulepszania modeli sztucznej inteligencji" (training and improving AI models), with the right to sublicense. Termination does not affect prior use.**
- Provider: **Fitatu S.A.**, ul. Wyspiańskiego 10/5, 60-749 Poznań, KRS 0000635344, NIP 7792444235, REGON 364839278 (§I.1).

**EU and Polish law context**
- CJEU C-30/14 *Ryanair v PR Aviation* (15 Jan 2015): where a database is protected by neither copyright nor the sui generis right, the Database Directive's user exceptions do not apply. The owner can restrict use, including screen scraping, by contract or terms of use. — [Kluwer Copyright Blog](https://legalblogs.wolterskluwer.com/copyright-blog/ryanair-ltd-v-pr-aviation-bv-contracts-rights-and-users-in-a-low-cost-database-law/); [IPKat](https://ipkitten.blogspot.com/2015/01/breaking-cjeu-says-that-owner-of-online.html); [IPPT judgment PDF](https://www.ippt.eu/sites/ippt/files/2015/IPPT20150115_ECJ_Ryanair_v_PR_Aviation.pdf)
- Database Directive 96/9/EC (implemented in Poland by the 2001 Act on protection of databases): sui generis right against extraction or re-utilisation of a substantial part of a database with substantial investment. — [EUR-Lex 31996L0009](https://eur-lex.europa.eu/legal-content/EN/ALL/?uri=CELEX:31996L0009)
- Software Directive 2009/24/EC, Art. 6:
  - Decompilation is allowed without authorisation only when indispensable to obtain information needed for interoperability of an independently created program, by a lawful user, when that information is not readily available, and confined to the necessary parts.
  - The information may not be used for other purposes or given to others except for that interoperability.
  - Art. 8: contractual terms contrary to Art. 6 are null.
  - — [EUR-Lex 32009L0024](https://eur-lex.europa.eu/legal-content/EN/ALL/?uri=CELEX:32009L0024)
- GDPR: Art. 9 special categories include health data, and processing needs e.g. explicit consent (Art. 9(2)(a)). Art. 20 gives a right to data portability for data the user provided, processed on consent or contract, in a structured, machine-readable format. — [EUR-Lex GDPR](https://eur-lex.europa.eu/eli/reg/2016/679/oj)
- The CJEU reads "health data" broadly. In C-21/23 *Lindenapotheke* (4 Oct 2024), data entered when ordering pharmacy-only medicines online were held to be health data. — [CURIA C-21/23](https://curia.europa.eu/juris/liste.jsf?num=C-21/23)

### Inferences
- **ToS breach is near-certain** for any server-side Coolinary integration via the private API. It is automated (§VI.4(a)), uses data for non-Fitatu purposes (§IX.3), and could be construed as third-party account use (§IV.6) even with the user's consent. Fitatu's remedies are contractual against the *user* (account blocking) and potentially civil claims (unfair competition, database right) against Coolinary. It can also break the integration at will by rotating client keys or adding attestation.
- **Database right:** Fitatu's dietitian-moderated database of over 1M items very likely involves "substantial investment". Bulk catalogue extraction would be high-risk. Reading or writing only the user's own diary entries is lower-risk under database law, but still a ToS breach. This is my assessment, not legal advice.
- **Directive 2009/24 Art. 6 does not help much.** It covers decompiling *software* to achieve interoperability. It does not authorise ongoing commercial use of a private server API against the operator's terms. The interface information is also arguably "readily available" via the open-source client, and that client was not built by Coolinary. Art. 8 may void §V.5's decompilation ban only in the narrow Art. 6 scope.
- **GDPR:** reading or writing a user's food diary, weight and goals would almost certainly be health data (Art. 9). Coolinary would need explicit consent, a DPIA, secure credential and token storage (Fitatu passwords or refresh tokens), and would hold another company's account credentials. The security, liability and reputational burden is significant for a solo founder.
- **Pushing content into Fitatu has its own cost.** Any recipe a Coolinary user saves into Fitatu, by any route, becomes "Materiały" that Fitatu may use, sublicense and use for AI training (§XIV.15). That gives Coolinary's recipe IP to a company running an overlapping feature.
- **Legitimate alternative for READ:** GDPR Art. 20 portability, i.e. the user's own XLS/CSV export (Premium) uploaded to Coolinary. It is lawful but manual and periodic, and gives no real-time "eaten today".

### Gaps
- I did not find any public case of Fitatu enforcing its ToS against third-party tools such as VikFit or fitatu-api, or any statement by Fitatu about them.
- I could not confirm whether the current Fitatu apps use device attestation or certificate pinning, which would affect how quickly unofficial clients break.
- No Polish court decision specific to API or ToS-based scraping of a nutrition app was found.

## 6. Company facts: ownership, users, countries/languages, pricing (PLN/EUR), database, AI, 2024–2026 news

### Takeaway
Fitatu is a profitable Poznań company founded by the Sympatia/Vitalia founders.
- **Ownership:** in early 2026 a 20-investor consortium (FIT SPV, organised by S20/SpeedUp) bought 53% for over 30M PLN, and installed CEO Jakub Biczkowski. It converted to **Fitatu S.A.** in 2026.
- **Users:** about 1.2M average MAU (peak about 1.5M in January 2026), about 240k paying subscribers, 5M+ Google Play installs.
- **Revenue:** over 80% from Poland. 2025 product revenue was about 20M PLN (+45% YoY).
- **Markets and product:** 7 languages; DACH is the priority market. It has a strong AI photo, voice and text logging push ("Fitatu AI 2.0 Acai").
- **Overlap with Coolinary:** Fitatu already has a Premium **"Lodówka/Smart Fridge"** (ingredients-to-recipes), meal plans and a shopping list.

### Cited Findings
**Ownership, legal and management**
- 3 Feb 2026: 20 investors bought 53% via SPV **FIT SPV**, in a deal "worth over 30 mln PLN" including a share buyback and capital injection. — [XYZ.pl](https://xyz.pl/tworcy-sympatii-oddali-kontrole-w-fitatu-grupa-znanych-inwestorow-rozkreci-zywieniowa-aplikacje-w-europie-zachodniej/); [MamBiznes, 3 Feb 2026](https://mambiznes.pl/news/nowy-rozdzial-aplikacji-fitatu-firma-zmienia-zarzad-i-planuje-globalna-ekspansje-narzedzia-do-liczenia-kalorii/)
  - Investors named: Maciej Zientara and Wojciech Duda (Supernova Group), Sławomir Nitek (ex-CEO Action Polska), André Gerstner, Ideo Software, Jakub Biczkowski, and S20/SpeedUp (Bartłomiej Gola).
  - New CEO: Jakub Biczkowski (ex-Wakacje.pl / Invia CEE).
  - Founders Wojciech Wrzaskała and Dominik Golenia (spelled "Goleń" by MamBiznes) keep about 26% (down from 72%) as advisors.
  - SpeedUp Innovation keeps its stake. It invested about 4M PLN in 2016 for about 18%.
  - Ewa Chodakowska and Lefteris Kavoukis sold out.
- Founders' history: Wrzaskała and Golenia sold Sympatia.pl to Onet in 2003. Fitatu grew out of Vitalia, whose diet business was sold in mid-2025. It pivoted to a mostly free calorie counter in 2019. — [XYZ.pl](https://xyz.pl/tworcy-sympatii-oddali-kontrole-w-fitatu-grupa-znanych-inwestorow-rozkreci-zywieniowa-aplikacje-w-europie-zachodniej/)
- Registry (aggregator, verify in KRS):
  - Legal form changed from Fitatu sp. z o.o. (KRS 0000635344) to **Fitatu S.A.**, new **KRS 0001266487**, published 22 Sep 2026. Share capital 101,920 PLN.
  - 10 Mar 2026 entry: CEO Jakub M. Biczkowski; supervisory board members Kuba Dudek and Marcin Jan Fejfer.
  - Shareholders named in filings include Fit SPV sp. z o.o., Fundacja Rodzinna Golenia, Fundacja Rodzinna Wrzaskałów and Speed Up Innovation sp. z o.o. S.K.A.
  - A 14 Oct 2026 entry under the *old* sp. z o.o. number lists Biczkowski and supervisory board members as "no longer" serving. This is probably the old register entry being closed after conversion: the new S.A. entry shows CEO "J. B." (inference).
  - — [NorthData](https://www.northdata.com/Fitatu+sp.+z+o.o.,+Pozna%C5%84/KRS0000635344)
- Conflict on legal name: the Terms "effective 13.07.2026" already name "Fitatu S.A." with the *old* KRS 0000635344, while the App Store still lists "Fitatu Sp. z o.o.". — [Fitatu Regulamin](https://www.fitatu.com/pl/regulations?p=959); [App Store PL](https://apps.apple.com/pl/app/fitatu-licznik-kalorii-i-dieta/id1011095795?l=pl)

**Financials**
- 2023 revenue 26.2M PLN, net profit about 3M PLN. 2024 revenue 29M PLN, profit about 2.4M PLN (both including Vitalia). 2025 Fitatu product only: about 20M PLN revenue (+45% YoY), EBITDA +130%. Over 90% of revenue is subscriptions and the rest ads. Over 80% of revenue is from Poland, with a target of about 60% within a few years. — [XYZ.pl](https://xyz.pl/tworcy-sympatii-oddali-kontrole-w-fitatu-grupa-znanych-inwestorow-rozkreci-zywieniowa-aplikacje-w-europie-zachodniej/)
- A search-result snippet citing KRS data gives revenue of 29,309,862 PLN and net profit of 3,311,551 PLN for an unspecified year, probably 2024. Unverified, and it differs from XYZ's ~2.4M profit. — [bizraport.pl (snippet only)](https://www.bizraport.pl/krs/0000635344/fitatu-spolka-z-ograniczona-odpowiedzialnoscia)

**Users (company-reported, unaudited)**

| Metric | Value | Source and date |
|---|---|---|
| Average MAU | ~1.2M | [XYZ](https://xyz.pl/tworcy-sympatii-oddali-kontrole-w-fitatu-grupa-znanych-inwestorow-rozkreci-zywieniowa-aplikacje-w-europie-zachodniej/) and [MamBiznes](https://mambiznes.pl/news/nowy-rozdzial-aplikacji-fitatu-firma-zmienia-zarzad-i-planuje-globalna-ekspansje-narzedzia-do-liczenia-kalorii/), Feb 2026 |
| Paying subscribers | ~240k | XYZ, Feb 2026 |
| Premium+AI users | ~150k | XYZ, Feb 2026 |
| Monthly active unique users | over 1.5M | [Spider's Web Bizblog, 14 Jan 2026](https://bizblog.spidersweb.pl/nowy-rok-nowy-polak-rekordowe-statystyki-fitatu) |
| DAU on 7 Jan 2026 | over 330k | Spider's Web, Jan 2026 |
| Downloads, first week of Jan 2026 | over 100k | Spider's Web, Jan 2026 |
| Meal photos analysed by AI, first week of Jan 2026 | 243k | Spider's Web, Jan 2026 |
| MAU in Poland | ~1M | [Cityboard, Aug 2024](https://cityboard.pl/blog/713-reklama-na-fitatu) |

  - Conflict: the 1.5M figure looks like a January peak, while 1.2M is an average.
- Google Play (fetched 30 Oct 2026): **5M+ downloads**; page metadata shows an install figure of 9,614,029. About 157K reviews. First released 23 Jun 2015, updated 14 Oct 2026. IAP "2,99 PLN – 169,99 PLN". The PL-store listing carries a "#1 top grossing health & fitness" badge. — [Google Play](https://play.google.com/store/apps/details?id=com.fitatu.tracker&hl=en_US)
- App Store PL: 4.8★ from about 101k ratings, version 4.17.3 (20 Oct). App Store DE: 4.8★ from 2,716 ratings. — [App Store PL](https://apps.apple.com/pl/app/fitatu-licznik-kalorii-i-dieta/id1011095795?l=pl); [App Store DE](https://apps.apple.com/de/app/id1011095795)

**Countries and languages**
- The app is in 7 languages (PL, EN, DE, ES, FR, IT, NL), and the website has /de/, /es/, /fr/, /it/ and /nl/ versions. — [App Store PL](https://apps.apple.com/pl/app/fitatu-licznik-kalorii-i-dieta/id1011095795?l=pl); [fitatu.com/pl](https://www.fitatu.com/pl)
- The press says "six languages", probably counting translations from Polish. — [MamBiznes](https://mambiznes.pl/news/nowy-rozdzial-aplikacji-fitatu-firma-zmienia-zarzad-i-planuje-globalna-ekspansje-narzedzia-do-liczenia-kalorii/)
- Largest user bases outside Poland: Germany, UK, France, Netherlands, Italy, Spain. DACH is the priority. The main competitor cited is Yazio. — [XYZ.pl](https://xyz.pl/tworcy-sympatii-oddali-kontrole-w-fitatu-grupa-znanych-inwestorow-rozkreci-zywieniowa-aplikacje-w-europie-zachodniej/)
- The API runs regional hosts such as `pl-pl.fitatu.com` and `de-de.fitatu.com`. — [GitHub fitatu_library](https://github.com/Jezue/fitatu_library)

**Pricing**
- PLN, App Store PL in-app purchase list, billing periods not shown: Premium 39,99 / 69,99 / 99,99 zł; Premium+AI 49,99 / 119,99 / 169,99 zł. — [App Store PL](https://apps.apple.com/pl/app/fitatu-licznik-kalorii-i-dieta/id1011095795?l=pl)
- A review states regular prices of Premium 39.99 zł/month and Premium+AI 49.99 zł/month, with frequent annual promos around 70–170 zł. — [Antyweb](https://antyweb.pl/subskrypcja-fitatu-premium)
- Premium+AI cost 169,99 zł/year in Mar 2025. — [Tabletowo (sponsored)](https://www.tabletowo.pl/fitatu-premium-ai-szacowanie-kalorii-ze-zdjec/)
- Promos seen: Premium+AI 135,99 zł/yr ([Pepper](https://www.pepper.pl/promocje/fitatu-premiumai-za-13599zl-w-subskrypcji-rocznej-1240563)), 99,99 zł/yr ([Pepper](https://www.pepper.pl/promocje/fitatu-premium-ai-9999-zl-za-rok-964971)), and Premium 74,99 zł/yr ([Pepper](https://www.pepper.pl/promocje/fitatu-premium-7499-zl-za-roczny-abonament-tluste-rabaty-832151)).
- Inference: 39,99/49,99 zł are monthly prices and 99,99/169,99 zł the standard annual prices.
- EUR, App Store DE, billing periods not shown: Premium €9.99, €19.99, €29.99, €34.99; Premium+AI €14.99, €47.90. — [App Store DE](https://apps.apple.com/de/app/id1011095795)
  - Inference: about €9.99/month for Premium and €14.99/month for Premium+AI. €47.90 is likely an annual or discounted Premium+AI plan.
- Fitatu says it is cheaper abroad but earns higher ARPU there than in Poland. — [XYZ.pl](https://xyz.pl/tworcy-sympatii-oddali-kontrole-w-fitatu-grupa-znanych-inwestorow-rozkreci-zywieniowa-aplikacje-w-europie-zachodniej/)
- The Terms say prices may differ by platform. — [Fitatu Regulamin §VII](https://www.fitatu.com/pl/regulations?p=959)

**Tiers (Terms §III.1)** — [Fitatu Regulamin](https://www.fitatu.com/pl/regulations?p=959)
- Free: tracking, about 1,000 recipes, Health Connect import.
- Premium: all recipes, "Wybór menu", **Lodówka (fridge)**, shopping list, "Uzupełnij dzienny cel" (complete daily goal), export, sending data to Health Connect and Apple Health, other fitness-app syncs, no ads.
- Premium+AI: everything in Premium, plus AI photo, voice and text estimation, each up to 600 uses per month.

**Database**
- "ponad milion produktów i potraw" (over a million products and dishes). — [Fitatu help: baza produktów](https://www.fitatu.com/pl/help/co-to-jest-baza-produktow-i-potraw-fitatu)
- The homepage claims coverage of 91% of products on the Polish market and about 100,000 new products a month. These are marketing claims; the second looks very high. — [fitatu.com/pl](https://www.fitatu.com/pl)
- Missing products can be added by AI label scanning. — [Fitatu help: baza produktów](https://www.fitatu.com/pl/help/co-to-jest-baza-produktow-i-potraw-fitatu)
- Tracks "39 vitamins and nutrients". — [Google Play](https://play.google.com/store/apps/details?id=com.fitatu.tracker&hl=en_US)
- Recipe count: "thousands" ([Google Play](https://play.google.com/store/apps/details?id=com.fitatu.tracker&hl=en_US)); "over 2000" for Premium in Mar 2025 ([Tabletowo](https://www.tabletowo.pl/fitatu-premium-ai-szacowanie-kalorii-ze-zdjec/)).

**AI and features (2025–2026)**
- "Fitatu AI 2.0 ACAI" photo estimation, AI voice search, AI text search and AI label scanning. — [App Store PL](https://apps.apple.com/pl/app/fitatu-licznik-kalorii-i-dieta/id1011095795?l=pl); [Google Play](https://play.google.com/store/apps/details?id=com.fitatu.tracker&hl=en_US)
- After an external study ranked it outside the top 5 for photo accuracy, Fitatu now claims to lead. — [XYZ.pl](https://xyz.pl/tworcy-sympatii-oddali-kontrole-w-fitatu-grupa-znanych-inwestorow-rozkreci-zywieniowa-aplikacje-w-europie-zachodniej/)
- AI photo feature was already live in Mar 2025. — [Tabletowo](https://www.tabletowo.pl/fitatu-premium-ai-szacowanie-kalorii-ze-zdjec/)
- **Smart Fridge / "Lodówka":** the user types available ingredients and gets catalogue recipes containing them. It is in the search screen's tabs (Szukaj / Własne / Ulubione / Lodówka) and is Premium. Nothing is said about calorie-goal scaling, whole-day planning or multiple family members. **[blog post 5 Jul 2022, pre-2024, but the feature is still advertised in the Oct 2026 Play listing and the 2026 Terms]** — [Fitatu Blog EN (2022)](https://www.fitatu.com/blog/en/not-sure-what-to-do-with-the-leftovers-in-your-fridge-take-advantage-of-the-fitatu-fridge-and-treat-yourself-to-a-meal/); [Google Play](https://play.google.com/store/apps/details?id=com.fitatu.tracker&hl=en_US); [Fitatu Regulamin §III.1](https://www.fitatu.com/pl/regulations?p=959)
- Other features: "Complete Daily Goal" smart suggestions, a shopping list auto-generated from planned meals, 7 meal plans, carbohydrate exchanges for diabetics, intermittent fasting, habits, and a web app at fitatu.com. — [Google Play](https://play.google.com/store/apps/details?id=com.fitatu.tracker&hl=en_US)

### Inferences
- Fitatu is a well-funded, profitable, growth-focused local leader with recent investor control and an international push. It has commercial reasons to protect its data moat and its Premium features (fridge, menu, shopping list, AI). That makes it unlikely to open diary access to an app with an adjacent "fridge-to-meals" proposition.
- Coolinary's differentiation versus Fitatu's Lodówka is household and multi-person planning, a whole day of meals (breakfast/dinner/supper), per-person portion scaling to calorie targets, and "not a diary". Fitatu's fridge is single-user, single-recipe search inside a diary app.
- Fitatu's paying base (about 240k) is roughly the population that could use a HealthKit/Health Connect READ route (Premium writes meals) or a "paste into text-AI" PUSH workaround (Premium+AI, about 150k).

### Gaps
- No official 2026 PLN or EUR price list with billing periods. The in-app purchase lists don't label periods.
- No audited MAU or download figures; all are company-reported.
- Product-database size is given only as "over a million". There is no exact recipe count.
- Employee count not found.

## 7. Evidence of recipe/meal-planning apps, recipe sites or caterers partnering with Fitatu

### Takeaway
I found **no evidence** that any Polish recipe site (Kwestia Smaku, Aniagotuje), meal-planning app or diet caterer has an official integration or content partnership with Fitatu. Brand partnerships are sponsored-recipe and co-marketing deals sold through Fitatu's marketing channel or agencies. Food-retail co-marketing exists (Lidl Plus, InPost, MediaMarkt), along with past celebrity-trainer ties (Ewa Chodakowska as a shareholder for diet-related cooperation, since exited).

### Cited Findings
- Searches combining Fitatu with Dietly, Kwestia Smaku or Aniagotuje returned no partnership material. — (search, 30 Oct 2026; e.g. [Fitatu help: Przepisy](https://www.fitatu.com/pl/help/co-znajde-w-%22przepisach%22) was the only Fitatu hit)
- Ewa Chodakowska and Lefteris Kavoukis had received Fitatu stakes "for diet-related cooperation" and sold them in the 2026 deal. — [XYZ.pl](https://xyz.pl/tworcy-sympatii-oddali-kontrole-w-fitatu-grupa-znanych-inwestorow-rozkreci-zywieniowa-aplikacje-w-europie-zachodniej/)
- The sponsored "recipe database integration" format puts brand recipes permanently in the catalogue. — [Cityboard, Aug 2024](https://cityboard.pl/blog/713-reklama-na-fitatu)
- Lidl Plus × Fitatu (Jan 2026) was a discount on Premium+AI, with no recipe integration. — [Fitatu Blog](https://www.fitatu.com/blog/nowy-rok-nowe-cele-zrob-pierwszy-krok-do-formy-z-fitatu-i-lidl-plus/)
- Seasonal own-content promos exist, such as the "Fit Święta z Fitatu®" holiday recipe e-book for Premium subscribers. — [Fitatu promo terms PDF](https://www.fitatu.com/blog/wp-content/uploads/Regulamin-promocji-%E2%80%9EFit-Swieta-z-Fitatu%C2%AE.pdf)
- Catering-to-Fitatu links exist only as unofficial tools (VikFit; fitatu-api `catering_sync.py`). — [vikfit.pl](https://vikfit.pl/); [GitHub PR #21](https://github.com/Jezue/fitatu_library/pull/21)

### Inferences
- There is no precedent for a third-party recipe or meal-planning app integrating with Fitatu. If Coolinary approached Fitatu it would be the first such case, and the deal Fitatu knows how to sell is "pay to place branded recipes or co-promote Premium".
- A co-marketing deal similar to Lidl Plus is conceivable once Coolinary has users. Coolinary could offer Fitatu Premium discounts to its users, or vice versa. This gives distribution, not data integration.

### Gaps
- Fitatu does not publish a partner list, so lesser-known recipe or caterer partnerships might exist without being searchable.

## 8. Realism of each integration route for Coolinary (synthesis)

### Takeaway
There is no official two-way integration. The only clean options are: (1) READ meal totals that Fitatu Premium users export to Apple Health or Health Connect, (2) deep links and copy-paste helpers into Fitatu's existing features, and (3) a paid or negotiated content and co-marketing deal. "Delegate the diary to Fitatu with one-tap logging" cannot currently be built legitimately. The unofficial API can do it technically, but breaches the ToS and is a GDPR health-data liability.

### Cited Findings
| Route | Direction | Feasibility today | Legal/ToS | Key evidence |
|---|---|---|---|---|
| Official API / partner API | R+W | None exists | n/a | [Fitatu współpraca](https://www.fitatu.com/pl/help/jakie-sa-mozliwosci-wspolpracy-z-fitatu) |
| Apple Health / Health Connect: read Fitatu's exported meals and weight | Read | Works only for Premium users who enabled "send data". Granularity unverified. | Clean (platform permissions + GDPR consent) | [Fitatu help](https://www.fitatu.com/en/help/can-fitatu-sync-with-other-devices-and-apps); [Regulamin §III.1](https://www.fitatu.com/pl/regulations?p=959) |
| Apple Health / Health Connect: Coolinary writes nutrition for Fitatu to pick up | Write | Fitatu documents reading only activity, burned kcal and weight, so likely not imported | Clean but likely useless for Fitatu | [Fitatu help PL](https://www.fitatu.com/pl/help/czy-jest-mozliwa-synchronizacja-fitatu-z-innymi-urzadzeniami-i-aplikacjami) |
| Deep links to Fitatu screens (`/app/search`, `/app/text-ai-estimation`, `/app/recipe/{id}`) | Write-assist | Works for opening existing content or screens; cannot create content | Clean | [AASA](https://www.fitatu.com/.well-known/apple-app-site-association) |
| "Copy for Fitatu" text, pasted into Fitatu text-AI | Write-assist | Premium+AI users only (~150k); 2–3 taps | Clean | [Regulamin §III.1](https://www.fitatu.com/pl/regulations?p=959); [XYZ](https://xyz.pl/tworcy-sympatii-oddali-kontrole-w-fitatu-grupa-znanych-inwestorow-rozkreci-zywieniowa-aplikacje-w-europie-zachodniej/) |
| Sponsored recipes in the Fitatu catalogue, plus deep links | Write | Paid and negotiated; static recipes, no per-person scaling | Clean; licence terms negotiated | [Cityboard](https://cityboard.pl/blog/713-reklama-na-fitatu) |
| Fitatu "share meal" links from a Coolinary account | Write | Untested; fixed portions; unknown link lifetime | Grey (built-in tool allowed by §V.4, but commercial/competitive clauses §XIV.9) | [Fitatu help: udostępnianie](https://www.fitatu.com/pl/help/czy-moge-udostepniac-posilki); [Regulamin](https://www.fitatu.com/pl/regulations?p=959) |
| User CSV/XLS export uploaded to Coolinary (GDPR Art. 20) | Read | Premium only; manual and periodic | Clean | [Regulamin §III.1](https://www.fitatu.com/pl/regulations?p=959); [GDPR](https://eur-lex.europa.eu/eli/reg/2016/679/oj) |
| Unofficial private API (fitatu-api), using the user's credentials | R+W | Technically full (diary, private recipes, goals, weight); breaks without notice | Breaches §VI.4(a), §IX.3, §IV.6; GDPR Art. 9 burden; credential custody | [GitHub](https://github.com/Jezue/fitatu_library); [Regulamin](https://www.fitatu.com/pl/regulations?p=959) |
| Co-marketing (Lidl Plus-style Premium discount swap) | Distribution | Plausible after traction | Clean | [Fitatu × Lidl Plus](https://www.fitatu.com/blog/nowy-rok-nowe-cele-zrob-pierwszy-krok-do-formy-z-fitatu-i-lidl-plus/) |

### Inferences
- If "calories already eaten today" matters, the most defensible MVP is the **HealthKit / Health Connect READ** route. It is generic, so it also works with other trackers that write nutrition, not only Fitatu. Fall back to a manual "I've eaten X kcal" input for non-Premium users.
- For the PUSH direction, the realistic MVP is a **"Log in Fitatu" helper**: show each person's portion as a gram-level ingredient list with a copy button, then deep-link to Fitatu search or text-AI. This keeps Coolinary "not a food diary" without promising an integration Fitatu does not offer.
- Approaching Fitatu (wspolpraca@fitatu.com) is worth it only as a low-cost exploratory e-mail. Present Coolinary as a household meal planner that sends users into Fitatu, not as a competitor to Lodówka or Wybór menu. Expect a marketing-deal answer, not API access.
- Building on the unofficial API would make a core product promise depend on a third party that can revoke it at will, and that contractually prohibits it.

### Gaps
- Hands-on verification is still needed: (a) the exact HealthKit and Health Connect data types and granularity Fitatu writes; (b) whether Fitatu ignores nutrition samples written by other apps; (c) share-link format, lifetime and behaviour for non-users; (d) whether deep links accept pre-fill parameters.
- Fitatu's appetite for a partnership with a startup is unknown and could only be learned by contacting them.
