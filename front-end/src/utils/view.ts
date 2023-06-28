/**
 * 色相環の最大角度
 */
const MAX_COLOR_WHEEL_ANGLE = 360;

/**
 * 彩度
 */
const SATURATION = "50%";

/**
 * 輝度
 */
const LIGHTNESS = "50%";

/**
 * 文字列から背景色のstyle属性値を生成する。
 * @param text 文字列
 * @returns style属性値
 */
export function createBgColorStyleFromText(text: string): string {
  return `background-color: ${createColorFromText(text)}`;
}

/**
 * 文字列をHSLカラーに変換する。
 * @param text 文字列
 * @returns HSLカラー
 */
function createColorFromText(text: string): string {
  // 1文字ずつ文字コードに変換して加算
  const textSum = Array.from(text)
    .map((val) => val.charCodeAt(0))
    .reduce((prev, cur) => prev + cur);

  const hue = textSum % MAX_COLOR_WHEEL_ANGLE;
  return `hsl(${hue}, ${SATURATION}, ${LIGHTNESS})`;
}
