import { MAGIC_NUMBERS } from "./consts";

export function IsNotAllowedIconFileFormat(dataUrl: string): boolean {
  // データURL(data:[<mediatype>][;base64],<data>)のデータ部をbase64デコード
  const byteString = atob(dataUrl.split(",")[1]);

  // バイト配列に変換
  const byteArray = new Uint8Array(byteString.length);
  for (let i = 0; i < byteString.length; i++) {
    byteArray[i] = byteString.charCodeAt(i);
  }

  for (const magicNum of Object.values(MAGIC_NUMBERS)) {
    const magicNumPartArray = byteArray.subarray(0, magicNum.length);
    if (magicNum.join() === magicNumPartArray.join()) {
      return false;
    }
  }

  return true;
}
