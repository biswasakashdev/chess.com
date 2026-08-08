export interface BoardTheme {
  name: string;
  lightSquare: string;
  darkSquare: string;
  selectedSquare: string;
  lastMoveHighlight?: string;
}

export const THEMES: Record<string, BoardTheme> = {
  classicWood: {
    name: "Classic Wood",
    lightSquare: "#f0d9b5",
    darkSquare: "#b58863",
    selectedSquare: "#baca44",
  },
  greenLichess: {
    name: "Lichess Green",
    lightSquare: "#eeeed2",
    darkSquare: "#769656",
    selectedSquare: "#bbcb2b",
  },
  oceanBlue: {
    name: "Ocean Blue",
    lightSquare: "#e1eec3",
    darkSquare: "#f05053",
    selectedSquare: "#7ea04d",
  },
  darkGlass: {
    name: "Dark Cyber",
    lightSquare: "#2a2e3d",
    darkSquare: "#1a1c23",
    selectedSquare: "#00adb5",
  },
};
