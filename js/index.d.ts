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
  extension?: string;
  key: string;
  label: string;
  help?: string;
  type: "checkbox" | "text" | "textarea" | "select" | "number";
  value?: unknown;
  options?: Option[];
  required?: boolean;
}

export interface Section {
  extension?: string;
  title: string;
  fields: Field[];
}

export interface ContributionRequest {
  entity: string;
  action: string;
  id?: string;
  record?: unknown;
  values?: Record<string, unknown>;
}

export interface ContributionResponse {
  sections: Section[];
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
  info: () => Info | Promise<Info>;
  contribute: (
    req: ContributionRequest,
  ) => ContributionResponse | Promise<ContributionResponse>;
  handle: (event: Event) => EventResponse | Promise<EventResponse>;
}

export function serve(extension: Extension): void;
