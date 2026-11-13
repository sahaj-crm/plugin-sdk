import {
  asBool,
  asString,
  extensionValues,
  serve,
} from "../../js/index.js";

const extensionID = "com.example.jsfull";

serve({
  contribute: () => ({
    fields: [{
      name: "Enable the full JavaScript feature",
      param: "feature_enabled",
      type: "checkbox",
      note: "Stored under the extension's own data.",
    }, {
      name: "Full JavaScript mode",
      param: "feature_mode",
      type: "select",
      note: "JavaScript select",
      options: [
        { value: "one", label: "1" },
        { value: "two", label: "2" },
      ],
    }],
  }),

  handle: (event) => {
    const input = event.input || {};
    const values = extensionValues(input.Extras, extensionID);
    console.error(
      `full js extension: ${event.phase} ${event.name} ` +
        `feature_enabled=${asBool(values.feature_enabled)} feature_mode=${asString(values.feature_mode)}`,
    );

    if (event.phase === "after") {
      const output = event.output || {};
      console.error(`full js extension: saved id=${output.ID}`);
    }

    return {};
  },
});
