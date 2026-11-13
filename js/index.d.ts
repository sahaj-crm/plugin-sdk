export type Phase = "before" | "after" | "failed";

export interface Info {
  id: string;
  name: string;
  version: string;
  api_version: number;
  targets?: string[];
  events?: string[];
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
  options?: Option[];
}

export interface ContributionRequest {
  entity: string;
  action: string;
  id?: string;
  record?: unknown;
  values?: Record<string, unknown>;
}

export interface ContributionResponse {
  fields: Field[];
}

export interface Event {
  name: string;
  phase: Phase;
  entity: string;
  action: string;
  id?: string;
  input?: unknown;
  output?: unknown;
  error?: string;
}

export interface EventResponse {
  error?: string;
  values?: Record<string, unknown>;
}

export interface Extension {
  contribute: (
    req: ContributionRequest,
  ) => ContributionResponse | Promise<ContributionResponse>;
  handle: (event: Event) => EventResponse | Promise<EventResponse>;
}

export function serve(extension: Extension): void;
