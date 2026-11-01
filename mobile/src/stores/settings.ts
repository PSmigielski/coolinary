import AsyncStorage from "@react-native-async-storage/async-storage";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

/** Coolinary shows kcal by default; kJ is an option. */
export type EnergyUnit = "kcal" | "kJ";

type SettingsState = {
  energyUnit: EnergyUnit;
  setEnergyUnit: (unit: EnergyUnit) => void;
};

/** Client-only preferences, persisted on the device (never synced to the server). */
export const useSettings = create<SettingsState>()(
  persist(
    (set) => ({
      energyUnit: "kcal",
      setEnergyUnit: (energyUnit) => set({ energyUnit }),
    }),
    {
      name: "coolinary.settings",
      version: 1,
      storage: createJSONStorage(() => AsyncStorage),
      partialize: ({ energyUnit }) => ({ energyUnit }),
    },
  ),
);
