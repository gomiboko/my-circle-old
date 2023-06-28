interface FormObject {
  [key: string]: string | File | null;
}

/**
 * Content-Type
 */
export const CONTENT_TYPE = Object.freeze({
  MULTIPART_FORM_DATA: Object.freeze({ "content-type": "multipart/form-data" }),
});

/**
 * FormDataを生成する。
 * 値がnullのものは除外される。
 * @param dataObj FormDataに設定するデータ
 * @returns FormDataオブジェクト
 */
export function createFormData(dataObj: FormObject): FormData {
  const data = new FormData();

  for (const key of Object.keys(dataObj)) {
    const val = dataObj[key];
    if (val !== null) {
      data.append(key, val);
    }
  }

  return data;
}
