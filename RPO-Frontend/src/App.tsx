import { FormEvent, useEffect, useMemo, useState } from "react";
import {
  Activity,
  CheckCircle2,
  CreditCard,
  KeyRound,
  LayoutDashboard,
  LogIn,
  LogOut,
  Pencil,
  Plus,
  ReceiptText,
  RefreshCw,
  Search,
  Server,
  ShieldCheck,
  Trash2,
  UserRound,
  WalletCards,
  XCircle
} from "lucide-react";
import { api, ApiRecord, clearStoredToken, getStoredToken, getValue, storeToken } from "./api";

type Route = "dashboard" | "terminals" | "cards" | "keys" | "transactions" | "users" | "authorize";
type FieldType = "text" | "number" | "money" | "toggle" | "select";

const POLLING_INTERVAL_MS = 5000;

interface FieldConfig {
  key: string;
  label: string;
  type: FieldType;
  accessors?: string[];
  create?: boolean;
  update?: boolean;
  required?: boolean;
  options?: Array<{ label: string; value: string | number }>;
}

interface ColumnConfig {
  label: string;
  accessors: string[];
  format?: (value: string | number | boolean | null | undefined, record: ApiRecord) => string;
}

interface ResourceConfig {
  route: Route;
  resource: string;
  title: string;
  description: string;
  icon: typeof Server;
  columns: ColumnConfig[];
  fields: FieldConfig[];
  canCreate?: boolean;
  canUpdate?: boolean;
  canDelete?: boolean;
  extraAction?: "balance";
}

const idColumn: ColumnConfig = { label: "№", accessors: ["ID", "id"] };

const resources: ResourceConfig[] = [
  {
    route: "terminals",
    resource: "terminals",
    title: "Терминалы",
    description: "Устройства приема операций и их адреса.",
    icon: Server,
    columns: [
      idColumn,
      { label: "Серийный номер", accessors: ["SerialNumber", "serialNumber"] },
      { label: "Название", accessors: ["Name", "name"] },
      { label: "Адрес", accessors: ["Address", "address"] }
    ],
    fields: [
      { key: "id", label: "№", type: "number", update: true },
      { key: "serial_number", label: "Серийный номер", type: "text", required: true, accessors: ["SerialNumber", "serialNumber"] },
      { key: "name", label: "Название", type: "text", required: true, accessors: ["Name", "name"] },
      { key: "address", label: "Адрес", type: "text", required: true, accessors: ["Address", "address"] }
    ],
    canCreate: true,
    canUpdate: true,
    canDelete: true
  },
  {
    route: "cards",
    resource: "cards",
    title: "Карты",
    description: "Баланс, владелец, блокировка и привязка ключа.",
    icon: CreditCard,
    columns: [
      idColumn,
      { label: "Номер", accessors: ["CardNumber", "cardNumber"] },
      { label: "Владелец", accessors: ["OwnerName", "ownerName"] },
      { label: "Баланс", accessors: ["Balance", "balance"], format: (value) => money(Number(value || 0)) },
      { label: "Статус", accessors: ["IsBlocked", "isBlocked"], format: (value) => (Number(value) === 1 ? "Заблокирована" : "Активна") },
      { label: "Ключ", accessors: ["KeyID", "keyID"] }
    ],
    fields: [
      { key: "id", label: "№", type: "number", update: true },
      { key: "card_number", label: "Номер карты", type: "text", required: true, accessors: ["CardNumber", "cardNumber"] },
      { key: "owner_name", label: "Владелец", type: "text", required: true, accessors: ["OwnerName", "ownerName"] },
      { key: "balance", label: "Баланс", type: "money", required: true, accessors: ["Balance", "balance"] },
      { key: "is_blocked", label: "Заблокирована", type: "toggle", accessors: ["IsBlocked", "isBlocked"] },
      { key: "key_id", label: "№ ключа", type: "number", accessors: ["KeyID", "keyID"] }
    ],
    canCreate: true,
    canUpdate: true,
    canDelete: true,
    extraAction: "balance"
  },
  {
    route: "keys",
    resource: "keys",
    title: "Ключи",
    description: "Значения ключей, используемых картами.",
    icon: KeyRound,
    columns: [idColumn, { label: "Значение", accessors: ["Value", "value"] }],
    fields: [
      { key: "id", label: "№", type: "number", update: true },
      { key: "value", label: "Значение", type: "text", required: true, accessors: ["Value", "value"] }
    ],
    canCreate: true,
    canUpdate: true,
    canDelete: true
  },
  {
    route: "transactions",
    resource: "transactions",
    title: "Транзакции",
    description: "История операций по картам и терминалам.",
    icon: ReceiptText,
    columns: [
      idColumn,
      { label: "Сумма", accessors: ["Amount", "amount"], format: (value) => money(Number(value || 0)) },
      { label: "Карта", accessors: ["CardID", "cardID"] },
      { label: "Терминал", accessors: ["TerminalID", "terminalID"] },
      { label: "Операция", accessors: ["Operation", "operation"], format: (value) => operationLabel(String(value || "")) },
      { label: "Создана", accessors: ["CreatedAt", "createdAt"], format: (value) => String(value || "").slice(0, 19).replace("T", " ") }
    ],
    fields: [
      { key: "id", label: "№", type: "number", update: true },
      { key: "amount", label: "Сумма", type: "money", required: true, accessors: ["Amount", "amount"] },
      { key: "card_id", label: "№ карты", type: "number", required: true, accessors: ["CardID", "cardID"] },
      { key: "terminal_id", label: "№ терминала", type: "number", required: true, accessors: ["TerminalID", "terminalID"] }
    ],
    canCreate: true,
    canDelete: true
  },
  {
    route: "users",
    resource: "users",
    title: "Пользователи",
    description: "Учетные записи и административный флаг.",
    icon: UserRound,
    columns: [
      idColumn,
      { label: "Логин", accessors: ["Login", "login"] },
      { label: "Имя", accessors: ["Name", "name"] },
      { label: "Админ", accessors: ["IsAdmin", "isAdmin"], format: (value) => (Number(value) === 1 ? "Да" : "Нет") },
      { label: "Хеш пароля", accessors: ["PasswordHash", "passwordHash"], format: (value) => shortHash(String(value || "")) }
    ],
    fields: [
      { key: "id", label: "№", type: "number", update: true },
      { key: "login", label: "Логин", type: "text", required: true, accessors: ["Login", "login"], create: true },
      { key: "name", label: "Имя", type: "text", required: true, accessors: ["Name", "name"] },
      { key: "password_hash", label: "Хеш пароля", type: "text", required: true, accessors: ["PasswordHash", "passwordHash"] },
      { key: "is_admin", label: "Администратор", type: "toggle", accessors: ["IsAdmin", "isAdmin"] }
    ],
    canCreate: true,
    canUpdate: true,
    canDelete: true
  }
];

const navItems = [
  { route: "dashboard" as const, label: "Обзор", icon: LayoutDashboard },
  ...resources.map((resource) => ({ route: resource.route, label: resource.title, icon: resource.icon })),
  { route: "authorize" as const, label: "Авторизация", icon: ShieldCheck }
];

function getRouteFromPath(): Route {
  const segment = window.location.pathname.split("/").filter(Boolean)[0] as Route | undefined;
  return navItems.some((item) => item.route === segment) ? segment! : "dashboard";
}

function money(value: number) {
  return new Intl.NumberFormat("ru-RU", { style: "currency", currency: "RUB" }).format(value);
}

function shortHash(value: string) {
  return value.length > 18 ? `${value.slice(0, 9)}...${value.slice(-6)}` : value;
}

function operationLabel(value: string) {
  if (value === "withdraw") return "Списание";
  if (value === "deposit") return "Пополнение";
  return value;
}

function recordId(record: ApiRecord) {
  return Number(getValue(record, ["ID", "id"], 0));
}

function normalizeFormValue(value: string, field: FieldConfig) {
  if (field.type === "toggle") {
    return value === "1" || value === "true" ? 1 : 0;
  }
  if (field.type === "number" || field.type === "money") {
    if (value === "" && field.key === "key_id") {
      return null;
    }
    return Number(value);
  }
  return value;
}

export function App() {
  const [route, setRoute] = useState<Route>(getRouteFromPath);
  const [token, setToken] = useState(getStoredToken);
  const [sessionMessage, setSessionMessage] = useState(token ? "Токен сохранен" : "Войдите для управления");

  useEffect(() => {
    const onPopState = () => setRoute(getRouteFromPath());
    window.addEventListener("popstate", onPopState);
    return () => window.removeEventListener("popstate", onPopState);
  }, []);

  useEffect(() => {
    if (!token) return;
    api
      .validate(token)
      .then((result) => setSessionMessage(result.message || "Токен активен"))
      .catch((error) => setSessionMessage(error.message));
  }, [token]);

  const navigate = (nextRoute: Route) => {
    setRoute(nextRoute);
    window.history.pushState(null, "", nextRoute === "dashboard" ? "/" : `/${nextRoute}`);
  };

  const logout = () => {
    clearStoredToken();
    setToken("");
    setSessionMessage("Сессия завершена");
  };

  if (!token) {
    return <LoginScreen onLogin={(nextToken) => {
      storeToken(nextToken);
      setToken(nextToken);
      setSessionMessage("Вход выполнен");
    }} />;
  }

  const currentResource = resources.find((resource) => resource.route === route);

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <div className="brand-mark">Р</div>
          <div>
            <strong>Панель РПО</strong>
            <span>Панель сервера</span>
          </div>
        </div>
        <nav className="nav-list">
          {navItems.map((item) => {
            const Icon = item.icon;
            return (
              <button
                key={item.route}
                className={route === item.route ? "nav-item active" : "nav-item"}
                onClick={() => navigate(item.route)}
              >
                <Icon size={18} />
                <span>{item.label}</span>
              </button>
            );
          })}
        </nav>
      </aside>

      <main className="workspace">
        <header className="topbar">
          <div>
            <span className="eyebrow">Адрес интерфейса: /api/v1</span>
            <h1>{route === "dashboard" ? "Обзор системы" : route === "authorize" ? "Авторизация операции" : currentResource?.title}</h1>
          </div>
          <div className="session">
            <span>{sessionMessage}</span>
            <button className="icon-button" onClick={logout} title="Выйти">
              <LogOut size={18} />
            </button>
          </div>
        </header>

        {route === "dashboard" && <Dashboard />}
        {route === "authorize" && <AuthorizePage />}
        {currentResource && <ResourcePage config={currentResource} token={token} />}
      </main>
    </div>
  );
}

function LoginScreen({ onLogin }: { onLogin: (token: string) => void }) {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setLoading(true);
    setMessage("");
    try {
      const result = await api.login(login, password);
      onLogin(result.token);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Не удалось войти");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <form className="login-card" onSubmit={submit}>
        <div className="login-icon">
          <LogIn size={28} />
        </div>
        <h1>Панель РПО</h1>
        <p>Войдите в учетную запись администратора, чтобы управлять API.</p>
        <label>
          Логин
          <input placeholder="Введите логин" value={login} onChange={(event) => setLogin(event.target.value)} required />
        </label>
        <label>
          Пароль
          <input type="password" placeholder="Введите пароль" value={password} onChange={(event) => setPassword(event.target.value)} required />
        </label>
        {message && <div className="alert error">{message}</div>}
        <button className="primary-button" disabled={loading}>
          {loading ? "Проверяем..." : "Войти"}
        </button>
      </form>
    </div>
  );
}

function Dashboard() {
  const [data, setData] = useState<Record<string, ApiRecord[]>>({});
  const [message, setMessage] = useState("Загрузка данных...");
  const [lastUpdated, setLastUpdated] = useState("");

  const load = async (silent = false) => {
    if (!silent) {
      setMessage("Загрузка данных...");
    }
    try {
      const entries = await Promise.all(resources.map(async (resource) => [resource.resource, await api.list(resource.resource)] as const));
      setData(Object.fromEntries(entries));
      const time = new Intl.DateTimeFormat("ru-RU", {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit"
      }).format(new Date());
      setLastUpdated(time);
      setMessage(silent ? `Автообновление: ${time}` : "Данные обновлены");
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Не удалось загрузить обзор");
    }
  };

  useEffect(() => {
    load();
    const timer = window.setInterval(() => load(true), POLLING_INTERVAL_MS);
    return () => window.clearInterval(timer);
  }, []);

  const cards = data.cards || [];
  const transactions = data.transactions || [];
  const terminals = data.terminals || [];
  const blockedCards = cards.filter((card) => Number(getValue(card, ["IsBlocked", "isBlocked"], 0)) === 1).length;
  const totalBalance = cards.reduce((sum, card) => sum + Number(getValue(card, ["Balance", "balance"], 0)), 0);

  return (
    <section className="content-grid">
      <div className="metric-card">
        <CreditCard size={22} />
        <span>Карты</span>
        <strong>{cards.length}</strong>
      </div>
      <div className="metric-card">
        <Server size={22} />
        <span>Терминалы</span>
        <strong>{terminals.length}</strong>
      </div>
      <div className="metric-card">
        <ReceiptText size={22} />
        <span>Транзакции</span>
        <strong>{transactions.length}</strong>
      </div>
      <div className="metric-card warning">
        <XCircle size={22} />
        <span>Блокировки</span>
        <strong>{blockedCards}</strong>
      </div>
      <div className="wide-panel">
        <div>
          <span className="eyebrow">Суммарный баланс</span>
          <h2>{money(totalBalance)}</h2>
        </div>
        <button className="secondary-button" onClick={() => load()}>
          <RefreshCw size={16} />
          Обновить
        </button>
      </div>
      <div className="wide-panel muted-panel">
        <Activity size={22} />
        <span>{message}{lastUpdated && !message.includes(lastUpdated) ? ` · Последнее обновление: ${lastUpdated}` : ""}</span>
      </div>
    </section>
  );
}

function ResourcePage({ config, token }: { config: ResourceConfig; token: string }) {
  const [rows, setRows] = useState<ApiRecord[]>([]);
  const [query, setQuery] = useState("");
  const [lookupId, setLookupId] = useState("");
  const [activeLookupId, setActiveLookupId] = useState("");
  const [message, setMessage] = useState("");
  const [modal, setModal] = useState<{ mode: "create" | "update" | "balance"; record?: ApiRecord } | null>(null);
  const [loading, setLoading] = useState(false);

  const load = async (silent = false) => {
    if (!silent) {
      setLoading(true);
      setMessage("");
    }
    try {
      setRows(await api.list(config.resource));
      if (!silent) {
        setActiveLookupId("");
      }
      if (!silent) {
        setMessage("Список обновлен");
      }
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Не удалось загрузить список");
    } finally {
      if (!silent) {
        setLoading(false);
      }
    }
  };

  useEffect(() => {
    load();
  }, [config.resource]);

  useEffect(() => {
    const timer = window.setInterval(() => {
      if (!document.hidden) {
        if (activeLookupId) {
          api
            .get(config.resource, Number(activeLookupId))
            .then((record) => setRows([record]))
            .catch((error) => setMessage(error instanceof Error ? error.message : "Не удалось обновить запись"));
        } else {
          load(true);
        }
      }
    }, POLLING_INTERVAL_MS);
    return () => window.clearInterval(timer);
  }, [config.resource, activeLookupId]);

  const filteredRows = useMemo(() => {
    const normalized = query.trim().toLowerCase();
    if (!normalized) return rows;
    return rows.filter((row) => JSON.stringify(row).toLowerCase().includes(normalized));
  }, [query, rows]);

  const findById = async () => {
    if (!lookupId) return;
    setLoading(true);
    try {
      setRows([await api.get(config.resource, Number(lookupId))]);
      setActiveLookupId(lookupId);
      setMessage(`Показана запись № ${lookupId}`);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Запись не найдена");
    } finally {
      setLoading(false);
    }
  };

  const remove = async (row: ApiRecord) => {
    const id = recordId(row);
    if (!window.confirm(`Удалить запись ID ${id}?`)) return;
    try {
      await api.delete(config.resource, id, token);
      setMessage("Запись удалена");
      await load();
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Не удалось удалить запись");
    }
  };

  return (
    <section className="resource-view">
      <div className="section-head">
        <div>
          <p>{config.description}</p>
        </div>
        <div className="toolbar">
          <div className="search-box">
            <Search size={16} />
            <input placeholder="Поиск по таблице" value={query} onChange={(event) => setQuery(event.target.value)} />
          </div>
          <div className="id-lookup">
            <input placeholder="№" value={lookupId} onChange={(event) => setLookupId(event.target.value)} />
            <button className="icon-button" onClick={findById} title="Найти по номеру">
              <Search size={17} />
            </button>
          </div>
          <button className="secondary-button" onClick={() => load()} disabled={loading}>
            <RefreshCw size={16} />
            Обновить
          </button>
          {config.canCreate && (
            <button className="primary-button" onClick={() => setModal({ mode: "create" })}>
              <Plus size={17} />
              Создать
            </button>
          )}
        </div>
      </div>

      {message && <div className="alert">{message}</div>}

      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              {config.columns.map((column) => (
                <th key={column.label}>{column.label}</th>
              ))}
              <th className="actions-col">Действия</th>
            </tr>
          </thead>
          <tbody>
            {filteredRows.map((row) => (
              <tr key={recordId(row)}>
                {config.columns.map((column) => {
                  const rawValue = getValue(row, column.accessors, "");
                  const value = column.format ? column.format(rawValue, row) : String(rawValue ?? "");
                  return <td key={column.label} data-label={column.label}>{value}</td>;
                })}
                <td className="row-actions" data-label="Действия">
                  {config.extraAction === "balance" && (
                    <button className="icon-button" onClick={() => setModal({ mode: "balance", record: row })} title="Изменить баланс">
                      <WalletCards size={17} />
                    </button>
                  )}
                  {config.canUpdate && (
                    <button className="icon-button" onClick={() => setModal({ mode: "update", record: row })} title="Изменить">
                      <Pencil size={17} />
                    </button>
                  )}
                  {config.canDelete && (
                    <button className="icon-button danger" onClick={() => remove(row)} title="Удалить">
                      <Trash2 size={17} />
                    </button>
                  )}
                </td>
              </tr>
            ))}
            {filteredRows.length === 0 && (
              <tr>
                <td colSpan={config.columns.length + 1} className="empty-cell">Нет данных</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {modal && (
        <EntityModal
          config={config}
          modal={modal}
          token={token}
          onClose={() => setModal(null)}
          onSaved={async (nextMessage) => {
            setMessage(nextMessage);
            setModal(null);
            await load();
          }}
        />
      )}
    </section>
  );
}

function EntityModal({
  config,
  modal,
  token,
  onClose,
  onSaved
}: {
  config: ResourceConfig;
  modal: { mode: "create" | "update" | "balance"; record?: ApiRecord };
  token: string;
  onClose: () => void;
  onSaved: (message: string) => void;
}) {
  const fields: FieldConfig[] = modal.mode === "balance"
    ? [
        { key: "id", label: "№", type: "number" as const, update: true },
        { key: "balance", label: "Новый баланс", type: "money" as const, required: true, accessors: ["Balance", "balance"] }
      ]
    : config.fields.filter((field) => {
        if (field.key === "id") return modal.mode === "update";
        if (modal.mode === "create") return field.create !== false && field.update !== true;
        return field.update !== false && field.create !== true;
      });

  const [form, setForm] = useState<Record<string, string>>(() => {
    const initial: Record<string, string> = {};
    const source = (modal.record || {}) as ApiRecord;
    for (const field of fields) {
      const value = field.key === "id"
        ? recordId(source)
        : getValue(source, field.accessors || [field.key], field.type === "toggle" ? 0 : "");
      initial[field.key] = String(value ?? "");
    }
    return initial;
  });
  const [message, setMessage] = useState("");

  const title = modal.mode === "create" ? `Создать: ${config.title}` : modal.mode === "balance" ? "Изменить баланс" : `Изменить: ${config.title}`;

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    const payload: ApiRecord = {};
    for (const field of fields) {
      payload[field.key] = normalizeFormValue(form[field.key] || "", field);
    }
    try {
      const result = modal.mode === "create"
        ? await api.create(config.resource, payload, token)
        : modal.mode === "balance"
          ? await api.updateCardBalance(payload, token)
          : await api.update(config.resource, payload, token);
      onSaved(result.message || "Сохранено");
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Не удалось сохранить");
    }
  };

  return (
    <div className="modal-backdrop" role="presentation" onMouseDown={onClose}>
      <form className="modal" onSubmit={submit} onMouseDown={(event) => event.stopPropagation()}>
        <div className="modal-head">
          <h2>{title}</h2>
          <button type="button" className="icon-button" onClick={onClose} title="Закрыть">×</button>
        </div>
        <div className="form-grid">
          {fields.map((field) => (
            <label key={field.key} className={field.type === "toggle" ? "toggle-row" : ""}>
              {field.label}
              {field.type === "toggle" ? (
                <select value={form[field.key] || "0"} onChange={(event) => setForm({ ...form, [field.key]: event.target.value })}>
                  <option value="0">Нет</option>
                  <option value="1">Да</option>
                </select>
              ) : field.type === "select" ? (
                <select value={form[field.key] || ""} onChange={(event) => setForm({ ...form, [field.key]: event.target.value })}>
                  {field.options?.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}
                </select>
              ) : (
                <input
                  type={field.type === "text" ? "text" : "number"}
                  step={field.type === "money" ? "0.01" : "1"}
                  value={form[field.key] || ""}
                  required={field.required}
                  readOnly={field.key === "id"}
                  onChange={(event) => setForm({ ...form, [field.key]: event.target.value })}
                />
              )}
            </label>
          ))}
        </div>
        {message && <div className="alert error">{message}</div>}
        <div className="modal-actions">
          <button type="button" className="secondary-button" onClick={onClose}>Отмена</button>
          <button className="primary-button">Сохранить</button>
        </div>
      </form>
    </div>
  );
}

function AuthorizePage() {
  const [form, setForm] = useState({
    card_number: "",
    terminal_serial_number: "",
    amount: "10",
    operation: "withdraw"
  });
  const [result, setResult] = useState<{ ok: boolean; text: string; balance?: number } | null>(null);
  const [loading, setLoading] = useState(false);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setLoading(true);
    try {
      const response = await api.authorize({
        card_number: form.card_number,
        terminal_serial_number: form.terminal_serial_number,
        amount: Number(form.amount),
        operation: form.operation
      });
      setResult({ ok: response.authorized, text: response.message, balance: response.balance });
    } catch (error) {
      setResult({ ok: false, text: error instanceof Error ? error.message : "Операция не выполнена" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="authorize-layout">
      <form className="operation-panel" onSubmit={submit}>
        <label>
          Номер карты
          <input placeholder="Введите номер карты" value={form.card_number} onChange={(event) => setForm({ ...form, card_number: event.target.value })} required />
        </label>
        <label>
          Серийный номер терминала
          <input placeholder="Введите серийный номер" value={form.terminal_serial_number} onChange={(event) => setForm({ ...form, terminal_serial_number: event.target.value })} required />
        </label>
        <label>
          Сумма
          <input type="number" step="0.01" min="0.01" value={form.amount} onChange={(event) => setForm({ ...form, amount: event.target.value })} required />
        </label>
        <label>
          Операция
          <select value={form.operation} onChange={(event) => setForm({ ...form, operation: event.target.value })}>
            <option value="withdraw">Списание</option>
            <option value="deposit">Пополнение</option>
          </select>
        </label>
        <button className="primary-button" disabled={loading}>
          <ShieldCheck size={17} />
          {loading ? "Выполняем..." : "Провести операцию"}
        </button>
      </form>
      <div className={result?.ok ? "result-panel success" : "result-panel"}>
        {result?.ok ? <CheckCircle2 size={34} /> : <XCircle size={34} />}
        <span>Результат операции</span>
        <h2>{result ? result.text : "Заполните форму и отправьте запрос"}</h2>
        {result?.balance !== undefined && <strong>{money(result.balance)}</strong>}
      </div>
    </section>
  );
}
