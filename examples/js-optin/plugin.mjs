import { serve } from "../../js/index.js";

serve({
  contribute: () => ({
    fields: [{
      name: "Enable the JavaScript feature",
      param: "js_flag",
      type: "checkbox",
      note: "Stored under the extension's own data.",
    }],
  }),

  handle: (event) => {
    console.error(`js extension: ${event.phase} ${event.name}`);
    return {};
  },
});
