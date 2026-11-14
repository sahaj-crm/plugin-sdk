export type Phase = "before" | "after" | "failed";

export interface Setting {
  key: string;
  name: string;
  type: string;
  mode?: string;
  placeholder?: string;
  note?: string;
  options?: Option[];
  value?: unknown;
  required?: boolean;
  secret?: boolean;
}

export interface Info {
  id: string;
  name: string;
  version: string;
  api_version: number;
  targets?: string[];
  events?: string[];
  settings?: Setting[];
}

export interface Option {
  value: string;
  label: string;
}

export interface Field {
  name?: string;
  param: string;
  value?: unknown;
  type: "checkbox" | "checkboxes" | "text" | "textarea" | "select" | "number";
  mode?: string;
  placeholder?: string;
  note?: string;
  options?: Option[];
}

export interface ContributionRequest {
  entity: string;
  action: string;
  id?: string;
  record?: unknown;
  values?: Record<string, unknown>;
  settings?: Record<string, unknown>;
}

export interface ContributionResponse {
  fields: Field[];
}

export type Extras = Record<string, Record<string, unknown>>;

export interface EntityInput<T = Record<string, unknown>> {
  Record: T;
  Extras: Extras;
}

export interface Event<TInput = unknown, TOutput = unknown> {
  name: string;
  phase: Phase;
  entity: string;
  action: string;
  id?: string;
  input?: TInput;
  output?: TOutput;
  error?: string;
  settings?: Record<string, unknown>;
}

export interface EventResponse {
  error?: string;
  patch?: Record<string, unknown>;
  values?: Record<string, unknown>;
}

export interface Extension {
  contribute: (
    req: ContributionRequest,
  ) => ContributionResponse | Promise<ContributionResponse>;
  handle: (event: Event) => EventResponse | Promise<EventResponse>;
}

export function serve(extension: Extension): void;

export function extensionValues(
  extras: Extras | undefined,
  extensionId: string,
): Record<string, unknown>;

export function asBool(value: unknown): boolean;
export function asString(value: unknown): string;
export function asStrings(value: unknown): string[];
