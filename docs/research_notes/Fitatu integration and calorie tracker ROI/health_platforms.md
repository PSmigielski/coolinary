# Apple HealthKit, Android Health Connect and Samsung Health Data SDK as a nutrition bridge between Coolinary and calorie trackers (state as of 30 Oct 2026)

Scope note: this covers platform mechanics and policies only. Each tracker's own API is covered by other researchers. "Confirmed" means read on a primary or official page (URL given). "Inference" means my reasoning and is flagged as such. Sources dated before 2024 are marked **[pre-2024, may be outdated]**.

---

## 1. HealthKit: which nutrition types exist, how food/meal metadata works, authorisation, cross-app reads, source visibility, iPad/Mac, 2025–2026 changes

### Takeaway
HealthKit has 39 dietary quantity types plus a `food` correlation that groups them into one food object. The only official food metadata key is `HKMetadataKeyFoodType`, a free-text "type of food" string. There is **no official meal-type key**. Users grant read and write separately per data type. Once read access is granted, an app can read samples written by other apps and see which app wrote each one (`sourceRevision`). However, it can never tell whether read access was denied, because denied data simply looks absent. HealthKit runs on iPhone and, since iPadOS 17, on iPad. On Mac the Health store cannot be read directly. I found no Apple-confirmed nutrition API changes for iOS 26 or 27.

### Cited Findings
- Apple's "Nutrition type identifiers" page lists 41 identifiers. Two are "essentials": `HKCorrelationTypeIdentifier.food` and `HKMetadataKeyFoodType`. The other 39 are quantity types:
  - **Macronutrients:** dietaryEnergyConsumed, dietaryCarbohydrates, dietaryFiber, dietarySugar, dietaryFatTotal, dietaryFatMonounsaturated, dietaryFatPolyunsaturated, dietaryFatSaturated, dietaryCholesterol, dietaryProtein
  - **Vitamins (13):** A, B6, B12, C, D, E, K, thiamin, riboflavin, niacin, pantothenic acid, biotin, folate
  - **Minerals (8):** calcium, chloride, iron, magnesium, phosphorus, potassium, sodium, zinc
  - **Ultratrace minerals (6):** chromium, copper, iodine, manganese, molybdenum, selenium
  - **Other:** dietaryWater, dietaryCaffeine
  
  — [Apple: Nutrition type identifiers](https://developer.apple.com/documentation/healthkit/nutrition-type-identifiers) (read via [JSON](https://developer.apple.com/tutorials/data/documentation/healthkit/nutrition-type-identifiers.json))
- Apple's guidance on that page:
  - Always provide the total (`dietaryFatTotal` or `dietaryCarbohydrates`) when recording subcategories. The subcategories should sum to no more than the total.
  - `dietaryEnergyConsumed` is measured in calories or kilojoules. Macronutrients are measured by mass.
  
  — [Apple: Nutrition type identifiers](https://developer.apple.com/tutorials/data/documentation/healthkit/nutrition-type-identifiers.json)
- `HKCorrelationTypeIdentifier.food` "combine[s] any number of nutritional samples into a single food object". Apple says: "When creating food samples, specify the type of food using the `HKMetadataKeyFoodType` metadata key", which is a String described as "The type of food that the HealthKit object represents". Availability: iOS 8+, iPadOS 8+, watchOS 2+, macOS 13+, visionOS 1+. — [Apple: HKCorrelationTypeIdentifier.food](https://developer.apple.com/documentation/healthkit/hkcorrelationtypeidentifier/food)
- Authorisation and privacy:
  - "Users can grant or deny permission separately for each type of data."
  - "an app isn't aware when the user denies permission to read data. From the app's point of view, no data of that type exists."
  - HealthKit data is stored locally on the device and encrypted while the device is locked. Background reads may therefore fail when the phone is locked, but writes are cached until unlock.
  - Apps must declare `NSHealthShareUsageDescription` (read) and `NSHealthUpdateUsageDescription` (write).
  
  — [Apple: Protecting user privacy (HealthKit)](https://developer.apple.com/documentation/healthkit/protecting-user-privacy)
- Source visibility: `HKObject.sourceRevision` is "The app or device that created this object". For app-saved samples it is set to the version of the app that saved the object. It is populated only on objects retrieved from the store. The older `HKObject.source` property is deprecated. — [Apple: HKObject.sourceRevision](https://developer.apple.com/documentation/healthkit/hkobject/sourcerevision)
- Reading other apps' nutrition works in practice. FoodNoms imports "nutrition entries that other apps save to Apple Health", and imported entries "appear in your Food Log alongside foods you log in Foodnoms". The same page says "Macs can't read Apple Health directly", so imports happen on the iPhone and then sync to the Mac. — [FoodNoms help: Reading from Health](https://foodnoms.com/help/reading-from-health)
- iPad: Apple announced at WWDC 2023 that the Health app would launch with iPadOS 17 and that "HealthKit is coming to iPad". Older documentation said HealthKit was unavailable on iPad; that is superseded. — [TechRadar on WWDC 2023](https://www.techradar.com/health-fitness/3-ways-apples-health-app-will-change-when-you-see-it-on-the-ipad-according-to-apple). A beta-era forum report said `isHealthDataAvailable()` returned false on an iPadOS 17 beta. — [Apple Developer Forums thread 732944](https://developer.apple.com/forums/thread/732944) **[2023, may be outdated]**
- 2025–2026 changes:
  - Bloomberg (Gurman) reported in 2025 that Apple was planning an integrated food diary in the Health app. It was first pegged to iOS 26.4. A 9to5Mac piece in April 2026 said bigger Health changes were still expected in iOS 27. — [9to5Mac, 30 Apr 2026](https://9to5mac.com/2026/04/30/ios-26-4-improved-apples-health-app-and-bigger-upgrades-are-coming-soon/); [T3](https://www.t3.com/tech/smartwatches/apples-health-app-could-be-getting-a-major-revamp-in-ios-26-4-with-food-tracking-and-an-ai-health-coach-rumoured)
  - Pre-release analysis (July 2026, vendor blog) says the redesigned iOS 27 Health app is "not a new API surface" and HealthKit read/write is "essentially unchanged". It also reports that iOS 27 food-photo analysis "does not sync to the Health app", while label scanning reportedly logs nutrients into Health. The author marks both claims "verify per beta"; neither is confirmed by Apple. — [Sahha blog, 7 Jul 2026](https://sahha.ai/blog/ios-27-health-fitness-developers/)

### Inferences
- **No official meal-type field.** Coolinary could encode meal type (breakfast/lunch/dinner) only as a custom metadata key, which other apps will not interpret. The food name goes in `HKMetadataKeyFoodType`. The Android side does have a `mealType` field (section 2), so the two platforms differ.
- **Recommended write format.** One `food` correlation per eaten portion, containing dietaryEnergyConsumed, protein, carbs, fat, fibre, sugar, saturated fat and sodium, with a time stamp equal to meal time. This mirrors what trackers write.
- **Updates and deletes.** Apple's `HKMetadataKeySyncIdentifier` / `HKMetadataKeySyncVersion` mechanism exists for this. I did not re-verify it in this session, so check it before implementation.
- **"Read today's intake" UX.** Because denial is invisible, Coolinary cannot tell "user denied read access" apart from "user logged nothing today". The UI should show "0 kcal found in Apple Health" rather than assume nothing was eaten. The feature should also be explicitly opt-in.
- **Double counting when summing.** Several trackers write the same day and the user may run two of them. Coolinary should either filter by `sourceRevision` (for example, let the user pick "my tracker") or use HKStatisticsQuery-style merged totals. I did not verify the exact deduplication semantics of statistics queries for dietary types here.
- **Locked-phone reads.** Reads only work when the phone is unlocked. That suits Coolinary, which would read when the user opens the app to plan the next meal.
- **Family scope.** HealthKit is one person's store on one device. Coolinary can only read and write data for the phone owner, not for other family members.

### Gaps
- I did not obtain Apple documentation that explicitly states "read permission for a type covers samples from all sources". This is inferred from FoodNoms, Lose It! and MacroFactor behaviour.
- I found no Apple-published iOS 26 or iOS 27 release notes that add or change dietary HealthKit types. Whether iOS 27's own food logging writes standard `dietary*` samples is unverified.

---

## 2. Health Connect (Android): NutritionRecord fields, MealType, permissions, history and background reads, platform status, EU availability

### Takeaway
`NutritionRecord` is rich. It has a time interval, energy, `energyFromFat`, about 10 macro fields, about 30 vitamin and mineral fields, a free-text `name` and an integer `mealType` (unknown, breakfast, lunch, dinner or snack). Access is gated by the `READ_NUTRITION` and `WRITE_NUTRITION` permissions. Reading other apps' data is limited by default to 30 days before the first grant; older data needs `READ_HEALTH_DATA_HISTORY`. Background reads need `READ_HEALTH_DATA_IN_BACKGROUND`. Health Connect is a built-in framework module from Android 14 and a Play Store app on Android 9–13. Each record carries the package name of the app that wrote it.

### Cited Findings
- `NutritionRecord` "Captures what nutrients were consumed as part of a meal or a food item."
  - **Required:** startTime, startZoneOffset, endTime, endZoneOffset, metadata.
  - **Optional energy (0–100,000 kcal):** energy, energyFromFat.
  - **Optional mass (0–100,000 g):** dietaryFiber, monounsaturatedFat, polyunsaturatedFat, protein, saturatedFat, sugar, totalCarbohydrate, totalFat, transFat, unsaturatedFat.
  - **Optional mass (0–100 g):** biotin, caffeine, calcium, chloride, cholesterol, chromium, copper, folate, folicAcid, iodine, iron, magnesium, manganese, molybdenum, niacin, pantothenicAcid, phosphorus, potassium, riboflavin, selenium, sodium, thiamin, vitamins A/B12/B6/C/D/E/K, zinc.
  - **Other optional:** `name` ("name of the food or drink, provided by the user") and `mealType` (default MEAL_TYPE_UNKNOWN).
  - **Validation:** startTime must be before endTime. On Android 14+ the platform performs the validation.
  
  — [AndroidX source: NutritionRecord.kt](https://raw.githubusercontent.com/androidx/androidx/androidx-main/health/connect/connect-client/src/main/java/androidx/health/connect/client/records/NutritionRecord.kt); [reference page](https://developer.android.com/reference/kotlin/androidx/health/connect/client/records/NutritionRecord)
- `MealType` constants: MEAL_TYPE_UNKNOWN = 0, MEAL_TYPE_BREAKFAST = 1, MEAL_TYPE_LUNCH = 2, MEAL_TYPE_DINNER = 3, MEAL_TYPE_SNACK = 4. — [AndroidX source: MealType.kt](https://raw.githubusercontent.com/androidx/androidx/androidx-main/health/connect/connect-client/src/main/java/androidx/health/connect/client/records/MealType.kt)
- Permissions are `android.permission.health.READ_NUTRITION` and `WRITE_NUTRITION` (prefix + "READ_NUTRITION"/"WRITE_NUTRITION"), plus `PERMISSION_READ_HEALTH_DATA_HISTORY` ("allows to read the entire history of health data") and `PERMISSION_READ_HEALTH_DATA_IN_BACKGROUND`. Both of the latter are "dependent on the version of HealthConnect installed" and must be checked via `HealthConnectFeatures`. — [AndroidX source: HealthPermission.kt](https://raw.githubusercontent.com/androidx/androidx/androidx-main/health/connect/connect-client/src/main/java/androidx/health/connect/client/permission/HealthPermission.kt)
- History limits:
  - "By default, all applications can read data from Health Connect for up to 30 days prior to when any permission was first granted."
  - On Android 14+ an app has no limit on reading its own data and a 30-day limit on other apps' data. On Android 13 and lower the 30-day limit applies to all data.
  - Without the history permission, reading older records "results in an error".
  
  — [Android Developers: Read data](https://developer.android.com/health-and-fitness/guides/health-connect/develop/read-data)
- Background reads and data origin:
  - Background reads require `READ_HEALTH_DATA_IN_BACKGROUND` and a feature check (`FEATURE_READ_HEALTH_DATA_IN_BACKGROUND`). Otherwise apps "normally read data from Health Connect when your app is in the foreground".
  - Reads can be filtered by writer via `dataOriginFilter = setOf(DataOrigin("<package>"))`.
  - Google recommends the aggregate API for totals because it "contains logic to handle duplicate records, and lessens the chances of rate limiting".
  
  — [Android Developers: Read data](https://developer.android.com/health-and-fitness/guides/health-connect/develop/read-data)
- Platform status and versions:
  - "Starting Android 14 (API Level 34), Health Connect is part of the Android Framework… there's no setup necessary" (Settings > Security and Privacy > Privacy Controls > Health Connect).
  - "On Android 13 (API Level 33) and lower… you need to install the Health Connect app from the Google Play Store."
  - The SDK supports Android 8+, but the Health Connect app requires Android 9+.
  - "Features tied to the system module remain unavailable on Android 13 and lower, even with the APK."
  - Apps check availability with `HealthConnectClient.getSdkStatus`.
  
  — [Android Developers: Get started](https://developer.android.com/health-and-fitness/guides/health-connect/develop/get-started)
- Google Fit is being retired: "Google Fit APIs will be supported until the end of 2026". The page was last updated 2026-03-10. — [Android Developers: Publish your health app](https://developer.android.com/health-and-fitness/health-connect/publish)

### Inferences
- **Android is the better write target.** Coolinary can write a `NutritionRecord` with `name` = dish name and `mealType` = breakfast/lunch/dinner/snack, which is semantically richer than HealthKit.
- **The 30-day limit does not matter for the core use case.** Reading "how much the user ate today" needs only today's data. `READ_HEALTH_DATA_HISTORY` should not be requested; it would add Play review friction for no benefit.
- **Skip the background permission.** Coolinary can read on app open, in the foreground, so `READ_HEALTH_DATA_IN_BACKGROUND` is unnecessary.
- **EU availability (inference, no explicit source found).** Health Connect is part of the Android framework from Android 14 and distributed via Play for Android 9–13. It should therefore be available on every Play-certified device in Poland and Germany. It would not be available on devices without Google Play services, such as recent Huawei phones.

### Gaps
- I found no official Google statement on Health Connect regional availability. EU availability is inferred.
- I did not verify any Android 16 or 2026 Health Connect changes specific to nutrition. One unrelated 2026 change was found: on-device steps now use a device-specific package name, per the read-data page.

---

## 3. Google Play Health Connect policy: declaration form, permitted uses, advertising ban, review timeline

### Takeaway
Every app that uses Health Connect must complete the **Health apps declaration form** (Play Console → App content). It must justify each requested data type with a user-facing feature and keep its privacy policy consistent across Play and Health Connect. "Providing personalized… meal plans" and "tracking nutrition" are explicitly permitted use cases. Selling health data and using it for any advertising, including personalised ads, are prohibited. Google publishes no review timeline; a vendor estimates up to 7 business days.

### Cited Findings
- Declaration form:
  - Health features are declared on the "Health apps" form on the App content page.
  - "For each permission requested, provide a clear and detailed justification explaining how your app uses the data to benefit the user."
  - "If your app does not require access to specific data types, you must not request access to them."
  - The process applies to new apps and to updates that change the set of data types. The form must be refilled if data-type needs change.
  - Approval is "allow-listed for a package name regardless of app version", but "your app's declaration may be reviewed again" on new versions.
  
  — [Android Developers: Publish your health app on Google Play](https://developer.android.com/health-and-fitness/health-connect/publish)
- Privacy policy and Data safety: the policy posted on the Play page "must be the same privacy policy that is displayed to users when they click the privacy policy link in Health Connect". The Data safety section must also be completed. — [Android Developers: Publish](https://developer.android.com/health-and-fitness/health-connect/publish). The rationale activity "must display the same privacy policy you provide for your app in the Google Play Console". — [Android Developers: Get started](https://developer.android.com/health-and-fitness/guides/health-connect/develop/get-started)
- Permitted use cases, under "Fitness, wellness, and coaching", include tracking nutrition and lifestyle habits and "providing personalized workout or meal plans". — [Google Play Console Help: Health permissions guidance](https://support.google.com/googleplay/android-developer/answer/12991134)
- Prohibited uses (same page):
  - Selling or transferring health data.
  - Using it for ads, "including personalized or interest-based advertising".
  - Using it for credit, insurance, employment or lending decisions.
  - "Sharing health data with third parties without explicit, informed user consent."
  
  The page also requires "Request the minimum data types needed and provide a valid use case for each request", and says sensitive categories (reproductive health, etc.) face "heightened scrutiny". — [Google Play Console Help: Health permissions guidance](https://support.google.com/googleplay/android-developer/answer/12991134)
- The same page lists `READ_HEALTH_DATA_IN_BACKGROUND` and `READ_HEALTH_DATA_HISTORY` as "other relevant permissions" but gives no separate justification rules. — [Google Play Console Help](https://support.google.com/googleplay/android-developer/answer/12991134)
- Review timeline: none stated by Google. A third-party SDK vendor says review "may take up to 7 business days". — [Sahha docs: Google Play Store submission](https://docs.sahha.ai/docs/guides/app-store-submission/google-play-store) (vendor estimate, not a Google commitment)
- Legacy-form deadline: apps that used the old form had to declare data types by 22 January 2025, or they could not submit changes for review. — [Android Developers: Publish](https://developer.android.com/health-and-fitness/health-connect/publish)

### Inferences
- **Coolinary fits the permitted use case cleanly** with READ_NUTRITION (to size the remaining meals) and WRITE_NUTRITION (to log meals marked as eaten). Requesting weight, steps or active energy should be avoided unless a feature truly uses them.
- **Ad-supported tiers.** If Coolinary ever runs ads, the ad stack must never receive Health-Connect-derived data. That includes the daily remaining-calorie figure if it was computed from Health Connect intake.

### Gaps
- No official review SLA. I could not find first-hand 2026 developer reports of typical approval times for nutrition-only declarations.

---

## 4. Samsung Health Data SDK (replacement for Samsung Health SDK for Android)

### Takeaway
The Samsung Health Data SDK (v1.1.0) supports both **reading and writing nutrition**. A released app, however, needs an approved **Samsung partnership**: package name and SHA-256 are registered, and writing needs a separate access code. Without approval the SDK works only in developer mode. The old Samsung Health SDK for Android was deprecated on 31 July 2025 and reaches end of service in 2028. Samsung Health also syncs nutrition with Health Connect, so for an indie app Health Connect covers Samsung users without any partnership.

### Cited Findings
- The Data SDK lists nutrition (and nutrition goal) as readable and nutrition as writable. It requires Samsung Health 6.30.2+ on Android 10+, does not support emulators, and is at version v1.1.0. — [Samsung: Health Data SDK overview](https://developer.samsung.com/health/data/overview.html)
- Process:
  - "Submit a partner request before releasing your app". After approval the package name and SHA-256 signature are registered. "Otherwise, the app using the Samsung Health Data SDK works only with the developer mode turned on."
  - Writing requires an access code issued after approval.
  - "The developer mode is intended only for testing or debugging your app."
  
  — [Samsung: Data SDK process](https://developer.samsung.com/health/data/process.html)
- FAQ: "To write health data into the Samsung Health, request partnership as the app creation process… After approval… an access code to write data on the developer mode will be provided. That access code is available for one application." — [Samsung: Data SDK FAQ](https://developer.samsung.com/health/data/faq.html)
- Deprecation: "Samsung Health SDK for Android was deprecated on July 31, 2025. The SDK will remain available for an additional two years and will reach the end of its service in 2028." — [Samsung Developer news, Oct 2025](https://developer.samsung.com/sdp/news/en/2025/10/30/dev-insight-oct-2025-move-to-samsung-health-data-sdk-as-samsung-health-sdk-for-android-deprecates-and-other-latest-news). The migration guide says writing to Samsung Health requires a partner request. — [Samsung: Migration guide](https://developer.samsung.com/health/data/migration-guide/overview.html)
- An older, undated Samsung overview page said Samsung "will not be accepting any applications for the Partner Apps Program at this time". Current Data SDK pages neither confirm nor deny this. — via search result for [Samsung Health Android SDK overview](https://developer.samsung.com/health/android/overview.html) **[undated, possibly outdated]**
- Samsung Health ↔ Health Connect: the blog's table maps "Nutrition" to `NutritionRecord`. "When Samsung Health has new or updated data, it writes the data to Health Connect… When Health Connect has updated data, Samsung Health retrieves it." The worked example only shows Samsung Health → Health Connect for nutrition. — [Samsung Developers blog, 18 Feb 2025](https://developer.samsung.com/health/blog/en/accessing-samsung-health-data-through-health-connect)

### Inferences
- **The Samsung SDK is not worth pursuing for Coolinary at launch.** It needs partnership approval of uncertain availability and duplicates what Health Connect already provides.
- **Samsung Health users are reachable via Health Connect,** at least in the Samsung Health → Health Connect direction. That direction is the one Coolinary needs for "read what I ate".

### Gaps
- Current (2026) partner-programme acceptance status and eligibility (individual vs company, EU) are not stated on the Data SDK pages.
- Whether Samsung Health's food diary actually shows nutrition records written by third-party apps into Health Connect is not confirmed by any official nutrition-specific statement. The two-way claim is generic.

---

## 5. Apple App Store Review Guidelines 5.1.3 / 5.1.2 and HealthKit rules

### Takeaway
HealthKit data must not be used for advertising, marketing or data mining, and must not be sold. It may be disclosed to third parties only with permission and only for health or fitness services. A privacy policy is mandatory. Apps must not write false or inaccurate data into HealthKit and **may not store personal health information in iCloud**. Sharing with third-party AI requires clear disclosure and explicit permission.

### Cited Findings
- 5.1.3(i): apps "may not use or disclose to third parties data gathered in the health, fitness, and medical research context—including from… HealthKit API… for advertising, marketing, or other use-based data mining purposes other than improving health management, or for the purpose of health research, and then only with permission… You must disclose the specific health data that you are collecting from the device." — [App Review Guidelines](https://developer.apple.com/app-store/review/guidelines/)
- 5.1.3(ii): "Apps must not write false or inaccurate data into HealthKit or any other medical research or health management apps, and may not store personal health information in iCloud." — [App Review Guidelines](https://developer.apple.com/app-store/review/guidelines/)
- 5.1.2(vi): data gathered from HealthKit "may not be used for marketing, advertising or use-based data mining, including by third parties". 5.1.2(i) also requires that apps "clearly disclose where personal data will be shared with third parties, including with third-party AI, and obtain explicit permission before doing so". — [App Review Guidelines](https://developer.apple.com/app-store/review/guidelines/)
- 2.5.1: "HealthKit should be used for health and fitness purposes and integrate with the Health app." — [App Review Guidelines](https://developer.apple.com/app-store/review/guidelines/)
- HealthKit developer rules:
  - "You must also provide a privacy policy for any app that uses the HealthKit framework."
  - "You must not disclose any information gained through HealthKit to a third party without express permission… Even with permission, you can only share… if they also provide a health or fitness service to the user."
  - "You can't sell information gained through HealthKit to advertising platforms, data brokers, or information resellers."
  - "Your app may not use information gained through… HealthKit for advertising or similar services."
  
  — [Apple: Protecting user privacy](https://developer.apple.com/documentation/healthkit/protecting-user-privacy)

### Inferences
- **"Inaccurate data" risk (important for Coolinary).** Writing *planned* meals, or meals computed for *other family members*, into the phone owner's HealthKit would arguably breach 5.1.3(ii). Coolinary should write only after an explicit "I ate this (my portion)" confirmation, using the user's own portion size, as Eat This Much does (section 7).
- **No iCloud storage.** If Coolinary syncs app state via CloudKit or iCloud, it must exclude HealthKit-derived values, such as today's intake read from Health.
- **Server-side AI meal generation.** If HealthKit-derived intake (for example "1,200 kcal eaten so far") is sent to Coolinary's server or an LLM provider to size the remaining meals, 5.1.2(i) requires explicit disclosure and permission. The HealthKit rules also limit third-party disclosure to recipients providing a health or fitness service. The lowest-risk design is to compute "remaining kcal" on device and send only the resulting budget. This matches the existing design, where only the daily calorie target goes to the server.

### Gaps
- The guidelines page did not show a "last updated" date. I did not check whether the 2026 revision added further health-specific wording beyond what is quoted.

---

## 6. Which calorie trackers WRITE nutrition to the OS stores, and which READ nutrition written by other apps

### Takeaway
**Writing** nutrition to Apple Health and Health Connect is standard: Fitatu, Yazio, MyFitnessPal, Lifesum, Cronometer, FatSecret, Lose It!, MacroFactor and Samsung Health (to Health Connect) all do it. **Reading foreign nutrition into the diary is rare.**
- **Fitatu, Yazio, MyFitnessPal, Lifesum, Cronometer and FatSecret do NOT** import food or nutrition logged by other apps. Their official help pages list only activity, weight, sleep and similar data as inbound.
- **Lose It!** does read it: foods from Apple Health appear in its log with a heart icon.
- **MacroFactor** imports daily *totals* only, and its own log overrides them.
- **FoodNoms** (an iOS niche app) imports entries.
- **Samsung Health:** unconfirmed for nutrition.

So **"Coolinary writes to the Health store → the meal appears in Fitatu, Yazio or MyFitnessPal" does not work.** The reverse direction does: Coolinary can read what any of these trackers wrote.

### Cited Findings
- **Fitatu.** Its integrations let you "send data on meals, weight, and height" and "receive data on burned calories, workouts, and activity". No import of other apps' meals is mentioned. It advises "To avoid duplicate data, it's best to use only one sync path." — [Fitatu help (EN)](https://www.fitatu.com/en/help/can-fitatu-sync-with-other-devices-and-apps); [Fitatu help (PL)](https://www.fitatu.com/pl/help/czy-jest-mozliwa-synchronizacja-fitatu-z-innymi-urzadzeniami-i-aplikacjami)
- **Yazio.**
  - Apple Health → Yazio covers "activities and workouts, body weight". Yazio → Apple Health covers "activity data, body values and nutrients". It is "not possible to connect YAZIO to more than one app at a time".
  - Source caveat: the article returned HTTP 403 to direct fetch, so this content comes from the search-engine rendering of the official page. — [YAZIO help: YAZIO and Apple Health](https://help.yazio.com/hc/en-us/articles/360018833998-YAZIO-and-Apple-Health)
  - A 2019 feature request said only macronutrients, not vitamins or minerals, were written to Apple Health. — [Yazio feature board](https://yazioen.featureupvote.com/suggestions/23881/sync-vitamins-minerals-with-apple-health) **[pre-2024, may be outdated]**
  - I found no Yazio help article covering nutrition sync to Health Connect. Its Google Fit → Health Connect article covers activity. — [YAZIO help: Syncing Google Fit with Health Connect](https://help.yazio.com/hc/en-us/articles/33129729688209-Syncing-Google-Fit-with-Health-Connect)
- **MyFitnessPal (Apple Health),** article updated 24 Sep 2026:
  - Food direction is MFP → Apple Health only: "Meal summaries and most nutrients sync. Caffeine and Vitamin A do not sync. Timestamps do not sync."
  - Exercise and sleep go Apple Health → MFP; weight syncs both ways.
  - Nothing is said about importing food from other apps.
  
  — [MFP help: Apple Health connection and syncing](https://support.myfitnesspal.com/hc/en-us/articles/360032271092-Apple-Health-connection-and-syncing)
- **MyFitnessPal (Health Connect),** article updated 10 Jul 2026:
  - MFP sends calories eaten "as a meal summary (not individual foods)". Only the default Breakfast/Lunch/Dinner/Snack headers map to meals; custom meals become "Other".
  - It receives cardio, glucose, sleep and steps. No inbound nutrition is listed.
  - Sync is "periodically throughout the day" or forced via Menu → Sync.
  
  — [MFP help: How To Link with Health Connect](https://support.myfitnesspal.com/hc/en-us/articles/10553948248973)
  - Conflict: a third-party blog claims MFP can import other apps' food via Dietary Energy with a heart icon. That wording matches Lose It!'s help article, not MFP's, and MFP's own help does not support it. — [amyfoodjournal blog](https://www.amyfoodjournal.com/blog/best-calorie-tracker-apple-health) (treated as unreliable)
- **Lifesum (Apple Health).** It exports food calories, food nutrition, weight and water. It imports exercise calories, steps, active calories, exercise time, workouts and weight. No food import. — [Lifesum help: Apple Health (NL)](https://help.lifesum.com/nl/article/how-to-connect-to-and-sync-with-apple-health-ios-wosoup)
- **Cronometer.**
  - It writes nutrition to Apple Health. A staff member said it "does not import any nutrient information from Apple Health", and another reply says "there isn't an option to read nutrition data from other sources (only workout data)". Users have requested this for years. — [Cronometer forum: import nutrition from Apple Health](https://forums.cronometer.com/discussion/5146/import-nutrition-from-apple-health-on-ios); [Cronometer forum: Import Food Data from Apple Health](https://forums.cronometer.com/discussion/6325/import-food-data-from-apple-health) (thread dates not confirmed; may be pre-2024)
  - Cronometer's blog lists inbound types as activity, exercise, height, sleep and heart rate, with no nutrition. — [Cronometer blog: Apple Health](https://cronometer.com/blog/apple-health)
- **FatSecret.** From Health it pulls calories burned and steps; it pushes food entries with nutrition details and syncs weight. No food import. — [FatSecret help: What data can Apple Health sync](https://www.fatsecret.com/fatsecret-app-help/third-party-apps/apple-health-data-sync-with-fatsecret)
- **Lose It!** (Apple Health article updated 1 Apr 2026):
  - It writes weight, food and exercise, and reads weight, food and exercise.
  - "Food from other third-party apps can appear… only if… permission is given to read Dietary Energy from Apple Health", and "Foods written from Apple Health to Lose It! will have a heart icon next to their listing in your Log."
  
  — [Lose It! help: Using Apple Health With Lose It (API copy)](https://loseit.zendesk.com/api/v2/help_center/en-us/articles/47773580085140.json); [article](https://loseit.zendesk.com/hc/en-us/articles/47773580085140-Using-Apple-Health-With-Lose-It)
  - Health Connect: nutrition, exercise, steps, weight and calories are listed as shared in both directions. The feature is "available to a small number of users as we gradually roll out". This comes from the search-engine rendering; direct fetch returned 403. — [Lose It! help: Using Health Connect With Lose It](https://loseit.zendesk.com/hc/en-us/articles/47650100219028-Using-Health-Connect-With-Lose-It)
- **MacroFactor.**
  - v2.9.3 (26 Aug 2024): it "can now import Calories, macros and micronutrients for the day from Apple Health and Health Connect", meaning daily totals rather than individual foods. — [MacroFactor 2.9.3 release notes](https://macrofactor.com/version-2-9-3/)
  - Priority order is manual estimate, then MacroFactor's food log, then Apple Health/Health Connect. If food is logged in MacroFactor that day, its entries override the external totals. — [MacroFactor help: Data sources](https://help.macrofactorapp.com/settings/data_sources/) (from search snippet)
  - A user reports that the import sometimes "doesn't seem to work". — [r/MacroFactor via redlib mirror](https://redlib.hbubli.cc/r/MacroFactor/comments/1kavrng/using_other_app_to_track_meals) (anecdotal)
- **FoodNoms** (iOS, paid tier): imports other apps' nutrition entries from Apple Health into the Food Log. Imported entries cannot be deleted in FoodNoms. — [FoodNoms help](https://foodnoms.com/help/reading-from-health)
- **Samsung Health:** writes nutrition to Health Connect, plus a generic two-way sync statement (section 4). A Samsung Community user reports that Samsung Health could not read a particular nutrition app's data. — [Samsung EU Community](https://eu.community.samsung.com/t5/mobile-apps-services/samsung-health-nutracheck/m-p/7457908/highlight/true) (anecdotal)
- **Downstream consumers also exist.** Wearable and coaching apps read tracker nutrition from Health Connect. For example, MyFitCoach says that with MFP connected to Health Connect "you no longer need to manually enter your calories and macronutrients". — [MyFitCoach help](https://intercom.help/myfitcoach/en/articles/11631892-activating-health-connect-synchronization-android). A Fitbit user reports Fitbit "is not reading nutrition from MFP" via Health Connect. — [Fitbit Community](https://community.fitbit.com/t5/Third-Party-Integrations/Do-we-need-to-move-to-Google-to-get-MFP-syncing-again/td-p/5644026) (anecdotal)

### Inferences
- **Decisive conclusion.** None of the three trackers the founder named (Fitatu, Yazio, MyFitnessPal) imports foreign nutrition from Apple Health or Health Connect, per their own help centres as of 2026. A "Coolinary writes → appears in Fitatu" flow is therefore **not achievable via the OS stores**. It would need each tracker's own partner API, which other researchers are covering. Of the major trackers, only Lose It! (strong in the US, weak in PL/DE) and MacroFactor (daily totals only) would reflect Coolinary-written meals.
- **The READ direction is viable and tracker-agnostic.** Fitatu (meals), Yazio (nutrients), MFP (meal summaries), Lifesum, Cronometer, FatSecret, Lose It! and MacroFactor all write nutrition. Coolinary can therefore read "kcal and macros eaten today" from HealthKit or Health Connect without integrating any tracker API. Caveats:
  - MFP does not sync timestamps to Apple Health, so meal timing is unreliable. Daily totals are fine.
  - Sync is periodic (MFP), so the value may lag by minutes to hours.
  - Users running two trackers will double count unless Coolinary filters by source.
- **Writing is still useful, for different reasons.** It helps Apple Health and Health Connect users who want a single record, and Lose It!/MacroFactor users. But if a user also logs the same meal manually in Fitatu, Lose It! or MacroFactor, the totals double count. Coolinary's own read logic must then exclude its own records, using `sourceRevision` or `DataOrigin`, when computing "eaten elsewhere".

### Gaps
- Yazio and Lose It! Health Connect details come from search-engine renderings because the official pages blocked direct fetch (HTTP 403).
- Cronometer forum statements are undated in my results and may be pre-2024.
- I found no official Samsung statement that Samsung Health's food diary displays third-party nutrition from Health Connect.
- I did not check the in-app permission sheets, which list the exact HealthKit types each tracker requests.

---

## 7. Real-world meal-planning or recipe apps that write meals to HealthKit or Health Connect

### Takeaway
The clearest precedent is **Eat This Much**. Its Apple Health sync is Premium-only and opt-in (disabled by default). It exports nutrition **only for meals the user marks "I ate this"** and syncs weight both ways. Mealime advertised Health integration but is being shut down on 21 October 2026. I found no evidence that Paprika, Samsung Food or Whisk write nutrition to the stores.

### Cited Findings
- Eat This Much (App Store listing):
  - "Enable Apple Health syncing for Premium users".
  - "nutrition information will automatically export for meals marked 'I ate this'".
  - Weight is imported and exported automatically.
  - Release note 2.0.160: Apple Health nutrition syncing is "opt-in only and disabled by default".
  
  — [App Store: Eat This Much](https://apps.apple.com/us/app/eat-this-much-meal-planner/id981637806)
- A user thread on Eat This Much syncing to Apple Health exists (Jan 2024). — [r/EatThisMuch via mirror](https://lr.us.psf.lt/r/EatThisMuch/comments/192ylxu/syncing_to_apple_health) (content not reviewed in detail)
- Mealime: "Mealime will be discontinued as of October 21, 2026". Albertsons acquired it in 2021; recipes move to Meals Hub, and "All personal data will be deleted". — [Mealime closing page](https://www.mealime.com/closing). Its App Store listing mentioned Health app integration, but the data written is not described. — [App Store: Mealime](https://apps.apple.com/US/app/id1079999103)
- Fitia, a vendor-style article, describes a meal plan that informs a calorie budget alongside Apple Health. — [Fitia blog](https://fitia.app/learn/article/meal-planning-app-weekly-grocery-lists-smart-devices/) (marketing source)

### Inferences
- **Eat This Much's "I ate this" + opt-in + per-user portion** pattern is the store-compliant template for Coolinary. It avoids the "inaccurate data" problem (5.1.3(ii)) and the "planned vs eaten" confusion.
- **Low adoption among planners** suggests the ROI of the write direction is modest for a planner. Most meal-planning apps do not bother; the user value lies in trackers.

### Gaps
- No reliable sources found on Paprika, Samsung Food or Whisk writing nutrition to HealthKit or Health Connect, and no user-experience studies of planners writing to Health.

---

## 8. GDPR: nutrition intake and weight as health data; consent for reading from a health store; does on-device-only processing reduce obligations?

### Takeaway
Weight tracked over time, and food or calorie tracking aimed at weight loss, are generally treated as **data concerning health** (GDPR Art. 4(15), Art. 9). Processing such data needs an Art. 9(2) exception, in practice **explicit consent**. Reading the data from the phone's health store also engages ePrivacy Art. 5(3) ("gaining access" to terminal equipment). There, the "strictly necessary for a service explicitly requested by the user" exemption likely covers a feature the user switches on. Keeping the data on the device strongly reduces exposure (no server-side Art. 9 processing, transfers or breach risk), but it does not automatically take the publisher out of GDPR. The household exemption shields the *user*, not the app maker.

### Cited Findings
- GDPR Art. 4(15) defines "data concerning health" as personal data related to physical or mental health "which reveal information about his or her health status". Art. 9(1) prohibits processing special categories unless an Art. 9(2) exception applies, such as (a) explicit consent. Art. 2(2)(c) excludes purely personal or household activity. — [GDPR (EUR-Lex)](https://eur-lex.europa.eu/eli/reg/2016/679/oj)
- Article 29 Working Party annex on health data in apps, 5 Feb 2015 (**pre-GDPR, under Directive 95/46; still widely cited but formally outdated**):
  - A single weight reading may not be health data, but "that aspect measured over time, especially in combination with age and sex, may be used to determine a significant aspect of an individual's health".
  - "When conclusions are drawn about someone's health, regardless of their reliability, these conclusions are to be treated as health data."
  - Calories-only diet apps that do not store foods are "unlikely" to yield health conclusions "unless the daily intake of calories is excessive in absolute terms". Combined with other data, health data "can be inferred".
  - For "apps to track food and exercise in an effort to lose weight", the "most likely derogation is explicit consent".
  - "If the data processing only takes place on the device itself, and no personal data are transmitted outside the device, the law wouldn't apply to the user, because of the exception for purely personal use".
  
  — [WP29 Annex – health data in apps and devices (PDF)](https://ec.europa.eu/justice/article-29/documentation/other-document/files/2015/20150205_letter_art29wp_ec_health_data_after_plenary_annex_en.pdf)
- EDPB Guidelines 2/2023 on the technical scope of Art. 5(3) ePrivacy Directive, final v2.0 adopted 7 Oct 2024. Art. 5(3) applies when there is "information", "terminal equipment of a subscriber or user" and "gaining access". It covers techniques beyond cookies. — [EDPB Guidelines 2/2023 v2 (PDF)](https://edpb.europa.eu/system/files/2024-10/edpb_guidelines_202302_technical_scope_art_53_eprivacydirective_v2_en_0.pdf); [Hunton summary](https://hunton.com/privacy-and-information-security-law/edpb-adopts-guidelines-on-scope-of-eprivacy-directive)
- CNIL mobile-app recommendations, adopted 18 Jul 2024, final published 2025. Apps "must obtain consent to process data that is not necessary for their operation". The recommendations explain how consent relates to OS technical permissions and list health data among the sensitive data apps access. I found no specific text on on-device-only processing. — [CNIL: Mobile applications recommendations](https://www.cnil.fr/en/mobile-applications-cnil-publishes-its-recommendations-better-privacy-protection)

### Inferences (legal reasoning, not legal advice; Polish UODO and German DSK/BfDI positions not checked)
- **Coolinary already processes health data.** Coolinary stores weight and goals on the phone and computes per-person calories. Food intake read from a health store, combined with a weight goal, falls squarely within the WP29 "track food… to lose weight" example, so it should be treated as Art. 9 data.
- **On-device-only processing.** If intake is read and used only on the device (never sent to Coolinary servers or analytics), the publisher arguably does not "process" it as controller in any server-side sense. This sharply cuts obligations: no Art. 9 storage, no DPA/processor contracts, no transfer issues, minimal breach exposure. But:
  - (a) the app code, which the publisher designs, still determines the purpose and means of access;
  - (b) ePrivacy Art. 5(3) covers reading from the device regardless of whether data leaves it;
  - (c) the WP29 domestic-exception sentence covers the *user*, not the publisher.
  
  Prudent practice is an explicit in-app opt-in screen naming the data (nutrition read and/or write) and its purpose, in addition to the OS permission dialog, and a privacy-policy section. Store rules (sections 3 and 5) require that privacy policy anyway.
- **The data that does leave the phone.** If the server-side "daily calorie target" is adjusted using intake read from Health (for example, "remaining kcal today"), that number is derived from health-store data. It is plausibly health data under the WP29 "conclusions" logic. Keeping the subtraction (target − eaten) on device and sending only "plan N kcal across these meals" is the cleaner design.
- **Children and family members.** Coolinary computes for other family members, but the health store holds only the phone owner's data. Writing other members' meals into the owner's store would be both inaccurate data (Apple 5.1.3(ii)) and misattributed personal data.

### Gaps
- No EDPB guidance found that is specific to reading HealthKit or Health Connect on device.
- CJEU case law on the broad reading of Art. 9 (e.g. C-184/20; C-21/23 on pharmacy orders, Oct 2024) was not verified in this session.
- No Polish UODO or German DPA guidance on fitness or nutrition apps was retrieved.

---

## 9. Engineering effort for a cross-platform app (React Native or Flutter plugins)

### Takeaway
Both ecosystems have maintained plugins with nutrition support:
- **Flutter:** `health` (carp.dk, v13.3.2) offers `writeMeal` and a `NUTRITION` type on both Apple Health and Health Connect, plus helpers for the Health Connect history and background permissions.
- **React Native:** `@kingstinct/react-native-healthkit` (Nitro modules, Expo config plugin, food correlation support) for iOS, and `react-native-health-connect` (Expo plugin built in since v4) for Android. The older `react-native-health` is rated inactive.

A read-today's-intake plus write-on-"I ate this" integration is a modest engineering task (inference: roughly 1.5–3 developer-weeks including both stores' compliance work).

### Cited Findings
- Flutter `health` 13.3.2 (verified publisher carp.dk; 682 likes; 160 pub points; Android and iOS):
  - Supports `writeMeal` on iOS and Android and a `NUTRITION` data type on both platforms, with `NutritionHealthValue` values.
  - Provides `isHealthDataHistoryAuthorized`/`requestHealthDataHistoryAuthorization` and `isHealthDataInBackgroundAvailable`/`requestHealthDataInBackgroundAuthorization`.
  - Google Fit support was removed in v11. iOS requires 15.0+. The device must be unlocked for iOS reads.
  
  — [pub.dev: health](https://pub.dev/packages/health)
- `@kingstinct/react-native-healthkit`:
  - "The library has been migrated to react-native-nitro-modules".
  - Expo is supported via a config plugin (custom dev client, not Expo Go).
  - Correlation types (food, blood pressure) can be queried, saved and subscribed to.
  - Metadata is supported, including custom string keys.
  
  — [README (jsDelivr)](https://cdn.jsdelivr.net/npm/@kingstinct/react-native-healthkit@latest/README.md); [npm](https://www.npmjs.com/package/%40kingstinct%2Freact-native-healthkit)
- `react-native-health-connect` (matinzd): `insertRecords` writes arrays of records in one transaction. The package "cannot be used in the Expo Go app" but has Expo integration built in since v4. — [react-native-health-connect docs: insertRecords](https://matinzd.github.io/react-native-health-connect/docs/api/methods/insertRecords)
- `react-native-health`: Snyk rates maintenance "inactive" (no commits in six months, no release in about a year; snapshot date not shown), with about 21.6k weekly downloads. — [Snyk: react-native-health](https://security.snyk.io/package/npm/react-native-health)
- Native Health Connect write example: Google's write-data guide shows a NutritionRecord example in which nutrients use Mass units and energy uses Energy units. — [Android Developers: Write data](https://developer.android.com/health-and-fitness/guides/health-connect/develop/write-data)

### Inferences (effort estimate, not sourced)
- **Feature 1, read today's intake.** Query today's dietary energy, protein, carbs and fat. On iOS use statistics or sample queries per type, excluding Coolinary's own `sourceRevision`. On Android use the aggregate API with an origin filter excluding Coolinary's package. Show a "from your tracker" figure in the planner.
  - Estimate: about 3–5 dev-days for both platforms with Flutter `health` or the RN libraries.
  - Plus about 2–3 days for edge cases: denied permissions being invisible on iOS, Health Connect not installed on Android 9–13, sync lag, double sources.
- **Feature 2, write eaten meal.** On iOS, a food correlation (or individual samples) with `HKMetadataKeyFoodType` = dish name. On Android, a `NutritionRecord` with `name` and `mealType`. Plus update or delete on undo.
  - Estimate: about 2–4 dev-days. Generic plugins may not expose the HealthKit correlation or sync-identifier fully, so a small native module may be needed.
- **Store compliance:**
  - Play Health apps declaration and privacy-policy alignment, the Health Connect rationale activity, Apple usage strings, App Privacy labels and a reviewer demo: about 2–4 days.
  - Plus Google review wait (no official SLA; vendor estimate up to 7 business days).
  - Plus possible App Review back-and-forth on health data.
- **Testing:**
  - Health Connect behaves differently on Android 13 (APK) and Android 14+ (framework module). Test both.
  - The Samsung Health Data SDK does not run on emulators (section 4), but Health Connect does.
- **Bridge fit.** The Flutter `health` plugin is the most "batteries included" single dependency for both stores. On React Native, the Kingstinct + matinzd pair is the current mainstream choice. Avoid `react-native-health` for new work.

### Gaps
- I did not verify the exact `NutritionHealthValue` / `writeMeal` field list in Flutter `health` 13.x, for example whether it writes HealthKit food correlations or separate samples, and whether it maps `mealType` and `name` on Health Connect. The package's API reference should be checked.
- I did not verify current versions or release dates of `react-native-health-connect` and `@kingstinct/react-native-healthkit`.
