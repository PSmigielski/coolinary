import type { EnergyUnit } from "@/stores/settings";

export const KJ_PER_KCAL = 4.184;

/** Formats an energy amount given in kcal in the user's chosen unit. */
export function formatEnergy(kcal: number, unit: EnergyUnit): string {
  const value = unit === "kJ" ? kcal * KJ_PER_KCAL : kcal;
  return `${Math.round(value).toLocaleString()} ${unit}`;
}
