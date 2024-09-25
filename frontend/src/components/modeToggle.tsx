import { Moon, Sun } from "lucide-react";
import { useEffect, useState } from "react";
import { useTheme } from "@/provider/themeProvider";
import { Switch } from "@/components/ui/switch.tsx";

export function ModeToggle() {
  const { theme, setTheme } = useTheme();
  const [isDarkMode, setIsDarkMode] = useState(theme === "dark");

  useEffect(() => {
    setIsDarkMode(theme === "dark");
  }, [theme]);

  const toggleTheme = () => {
    const newTheme = isDarkMode ? "light" : "dark";
    setTheme(newTheme);
  };

  return (
    <div className="flex items-center">
      <Sun className="h-[1.2rem] w-[1.2rem] mr-2 dark:text-white text-black" />
      <Switch checked={isDarkMode} onCheckedChange={toggleTheme} />
      <Moon className="h-[1.2rem] w-[1.2rem] ml-2  dark:text-white text-black" />
      <span className="sr-only">Toggle theme</span>
    </div>
  );
}
