import { asBool, pluginValues, serve } from "../../js/index.js";

const pluginID = "com.example.jssimple";

serve({
  contribute: () => ({
    fields: [{
      name: "Enable the simple JavaScript feature",
      key: "enabled",
      type: "checkbox",
      note: "Minimal example with one checkbox.",
    }],
  }),

  handle: (event) => {
    const input = event.input || {};
    const values = pluginValues(input.Extras, pluginID);
    console.error(
      `simple js plugin: ${event.phase} ${event.name} ` +
        `enabled=${asBool(values.enabled)}`,
    );
    return {};
  },
});
