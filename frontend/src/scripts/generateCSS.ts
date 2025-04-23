// generateCSS.ts
import { writeFileSync } from 'fs';
import { PRIMARY_COLORS, SECONDARY_COLORS } from '../styles/colors.ts';

const generateCSS = () => {
  let cssContent = ':root {\n';

  // PRIMARY_COLORSをCSS変数に変換
  Object.entries(PRIMARY_COLORS).forEach(([key, value]) => {
    cssContent += `  --primary-${key}: ${value};\n`;
  });

  // SECONDARY_COLORSをCSS変数に変換
  Object.entries(SECONDARY_COLORS).forEach(([key, value]) => {
    cssContent += `  --secondary-${key}: ${value};\n`;
  });

  cssContent += '}';

  // 変換結果をvariables.cssに出力
  writeFileSync('src/styles/variables.css', cssContent);
};

generateCSS();
