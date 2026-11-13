import { asBool, extensionValues, serve } from "../../js/index.js";

const extensionID = "com.example.jssimple";

serve({
  contribute: () => ({
    fields: [{
      name: "Enable the simple JavaScript feature",
      param: "enabled",
      type: "checkbox",
      note: "Minimal example with one checkbox.",
    }],
  }),

  handle: (event) => {
    const input = event.input || {};
    const values = extensionValues(input.Extras, extensionID);
    console.error(
      `simple js extension: ${event.phase} ${event.name} ` +
        `enabled=${asBool(values.enabled)}`,
    );
    return {};
  },
});
