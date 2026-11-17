import {
  asBool,
  asString,
  extensionValues,
  serve,
} from "../../js/index.js";

const extensionID = "com.example.jsfull";

serve({
  contribute: (req) => {
    const settings = req.settings || {};
    const defaultMode = asString(settings.default_mode) || "one";
    return {
      fields: [{
        name: "Enable the full JavaScript feature",
        key: "feature_enabled",
        type: "checkbox",
        note: "Stored under the extension's own data.",
      }, {
        name: "Full JavaScript mode",
        key: "feature_mode",
        type: "select",
        value: defaultMode,
        note: "Defaults from the extension setting.",
        options: [
          { value: "one", label: "Mode 1" },
          { value: "two", label: "Mode 2" },
        ],
      }],
    };
  },

  handle: (event) => {
    const input = event.input || {};
    const values = extensionValues(input.Extras, extensionID);
    const settings = event.settings || {};
    console.error(
      `full js extension: ${event.phase} ${event.name} ` +
        `feature_enabled=${asBool(values.feature_enabled)} feature_mode=${asString(values.feature_mode)} ` +
        `settings.default_mode=${asString(settings.default_mode)} ` +
        `settings.api_key_configured=${asString(settings.api_key) !== ""}`,
    );

    if (event.phase === "after") {
      const output = event.output || {};
      console.error(`full js extension: saved id=${output.ID}`);
    }

    return {};
  },
});
