import { serve } from "../../js/index.js";

serve({
  contribute: () => ({
    sections: [{
      title: "JavaScript Options",
      fields: [{
        key: "js_flag",
        label: "Enable the JavaScript feature",
        type: "checkbox",
      }],
    }],
  }),

  handle: (event) => {
    console.error(`js extension: ${event.phase} ${event.name}`);
    return {};
  },
});
