const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "/api/v1";
const TOKEN_KEY = "rpo.jwt";

export type ApiRecord = Record<string, string | number | boolean | null | undefined>;

export interface LoginResponse {
  token: string;
  message: string;
  user_id: number;
}

export interface AuthorizationResponse {
  authorized: boolean;
  message: string;
  card_number?: string;
  operation?: string;
  balance?: number;
}

export function getStoredToken() {
  return localStorage.getItem(TOKEN_KEY) || "";
}

export function storeToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearStoredToken() {
  localStorage.removeItem(TOKEN_KEY);
}

export function getValue<T = unknown>(record: ApiRecord, keys: string[], fallback: T): T {
  for (const key of keys) {
    const value = record[key];
    if (value !== undefined && value !== null) {
      return value as T;
    }
  }
  return fallback;
}

function translateApiMessage(message: string) {
  const clean = message.trim();
  const exact: Record<string, string> = {
    "Authorization successful": "Авторизация выполнена",
    "Token is valid": "Токен действителен",
    "Invalid or expired token": "Токен недействителен или истек",
    "Invalid request body": "Некорректное тело запроса",
    "Login and password are required": "Введите логин и пароль",
    "Invalid login or password": "Неверный логин или пароль",
    "Method not allowed": "Метод запроса не поддерживается",
    "authorization header is required": "Требуется заголовок авторизации",
    "invalid authorization header format": "Некорректный формат заголовка авторизации",
    "token is required": "Требуется токен",
    "forbidden: admin access required": "Требуются права администратора",
    "Card created": "Карта создана",
    "Card updated": "Карта обновлена",
    "Card balance updated": "Баланс карты обновлен",
    "Card deleted": "Карта удалена",
    "Key created": "Ключ создан",
    "Key updated": "Ключ обновлен",
    "Key deleted": "Ключ удален",
    "Terminal created": "Терминал создан",
    "Terminal updated": "Терминал обновлен",
    "Terminal deleted": "Терминал удален",
    "Transaction created": "Транзакция создана",
    "Transaction deleted": "Транзакция удалена",
    "User created": "Пользователь создан",
    "User updated": "Пользователь обновлен",
    "User deleted": "Пользователь удален",
    "Transaction processed successfully": "Операция выполнена",
    "Card number cannot be empty": "Введите номер карты",
    "Terminal serial number cannot be empty": "Введите серийный номер терминала",
    "Transaction amount must be greater than 0": "Сумма операции должна быть больше нуля",
    "Operation must be withdraw or deposit": "Выберите списание или пополнение",
    "card is blocked": "Карта заблокирована"
  };

  if (exact[clean]) {
    return exact[clean];
  }

  return clean
    .replace("invalid token:", "Некорректный токен:")
    .replace("failed to parse token:", "Не удалось прочитать токен:")
    .replace("failed to create card:", "Не удалось создать карту:")
    .replace("failed to update card", "Не удалось обновить карту")
    .replace("failed to delete card", "Не удалось удалить карту")
    .replace("failed to create terminal:", "Не удалось создать терминал:")
    .replace("failed to update terminal", "Не удалось обновить терминал")
    .replace("failed to delete terminal", "Не удалось удалить терминал")
    .replace("failed to create key:", "Не удалось создать ключ:")
    .replace("failed to update key", "Не удалось обновить ключ")
    .replace("failed to delete key", "Не удалось удалить ключ")
    .replace("failed to create transaction:", "Не удалось создать транзакцию:")
    .replace("failed to delete transaction", "Не удалось удалить транзакцию")
    .replace("failed to retrieve", "Не удалось получить")
    .replace("error finding card by number:", "Не удалось найти карту по номеру:")
    .replace("error finding terminal by serial number:", "Не удалось найти терминал по серийному номеру:")
    .replace("insufficient funds.", "Недостаточно средств.")
    .replace("required:", "требуется:")
    .replace("available:", "доступно:")
    .replace("card number cannot be empty", "номер карты не может быть пустым")
    .replace("owner name cannot be empty", "имя владельца не может быть пустым")
    .replace("key value cannot be empty", "значение ключа не может быть пустым")
    .replace("serial number cannot be empty", "серийный номер не может быть пустым")
    .replace("address cannot be empty", "адрес не может быть пустым")
    .replace("terminal name cannot be empty", "название терминала не может быть пустым")
    .replace("login cannot be empty", "логин не может быть пустым")
    .replace("name cannot be empty", "имя не может быть пустым")
    .replace("password hash cannot be empty", "хеш пароля не может быть пустым");
}

function translatePayload<T>(payload: T): T {
  if (payload && typeof payload === "object" && "message" in payload) {
    return {
      ...payload,
      message: translateApiMessage(String((payload as { message?: unknown }).message || ""))
    };
  }
  return payload;
}

async function request<T>(
  path: string,
  options: {
    method?: string;
    body?: unknown;
    token?: string;
  } = {}
): Promise<T> {
  const headers: Record<string, string> = {};

  if (options.body !== undefined) {
    headers["Content-Type"] = "application/json";
  }
  if (options.token) {
    headers.Authorization = `Bearer ${options.token}`;
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: options.method || "GET",
    headers,
    body: options.body === undefined ? undefined : JSON.stringify(options.body)
  });

  const contentType = response.headers.get("content-type") || "";
  const payload = contentType.includes("application/json")
    ? await response.json()
    : await response.text();

  if (!response.ok) {
    const message =
      typeof payload === "string" ? payload.trim() : payload?.message || "Ошибка запроса";
    throw new Error(translateApiMessage(message || `Код ответа ${response.status}`));
  }

  return translatePayload(payload as T);
}

export const api = {
  login: (login: string, password: string) =>
    request<LoginResponse>("/auth/login", {
      method: "POST",
      body: { login, password }
    }),
  validate: (token: string) =>
    request<{ valid: boolean; message: string }>("/auth/validate", {
      method: "POST",
      token
    }),
  list: (resource: string) => request<ApiRecord[]>(`/${resource}/all`),
  get: (resource: string, id: number) => request<ApiRecord>(`/${resource}/get?id=${id}`),
  create: (resource: string, body: ApiRecord, token: string) =>
    request<{ message: string }>(`/${resource}/create`, {
      method: "POST",
      body,
      token
    }),
  update: (resource: string, body: ApiRecord, token: string) =>
    request<{ message: string }>(`/${resource}/update`, {
      method: "PUT",
      body,
      token
    }),
  updateCardBalance: (body: ApiRecord, token: string) =>
    request<{ message: string }>("/cards/balance", {
      method: "PUT",
      body,
      token
    }),
  delete: (resource: string, id: number, token: string) =>
    request<{ message: string }>(`/${resource}/delete?id=${id}`, {
      method: "DELETE",
      token
    }),
  authorize: (body: ApiRecord) =>
    request<AuthorizationResponse>("/transactions/authorize", {
      method: "POST",
      body
    })
};
