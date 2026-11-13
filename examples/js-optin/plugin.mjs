import { serve } from "../../js/index.js";

serve({
  contribute: () => ({
    fields: [{
      name: "Enable the JavaScript feature",
      param: "ext[com.example.jsoptin][js_flag]",
      type: "checkbox",
    }],
  }),

  handle: (event) => {
    console.error(`js extension: ${event.phase} ${event.name}`);
    return {};
  },
});
