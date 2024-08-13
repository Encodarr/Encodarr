import dark from "./dark";
import light from "./light";

const defaultDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
const auto = defaultDark ? { ...dark } : { ...light };

const themes = {
  auto: auto,
  light: light,
  dark: dark,
};
export function getTheme(themeName: "auto" | "light" | "dark") {
  return themes[themeName];
}
export default themes;
