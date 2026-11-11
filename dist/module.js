"use strict";
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

// src/module.ts
var module_exports = {};
__export(module_exports, {
  plugin: () => plugin
});
module.exports = __toCommonJS(module_exports);
var import_data2 = require("@grafana/data");

// src/DataSource.ts
var import_data = require("@grafana/data");
var DataSource = class extends import_data.DataSourceApi {
  constructor(instanceSettings) {
    super(instanceSettings);
  }
  async query(request) {
    return { data: [] };
  }
  async testDatasource() {
    return { status: "success", message: "Ping datasource configured" };
  }
};

// src/components/ConfigEditor.tsx
var import_ui = require("@grafana/ui");
var import_jsx_runtime = require("react/jsx-runtime");
function ConfigEditor({ options, onOptionsChange }) {
  const jsonData = options.jsonData || {};
  const onIntervalChange = (event) => {
    onOptionsChange({
      ...options,
      jsonData: { ...jsonData, defaultInterval: event.target.value }
    });
  };
  const onTimeoutChange = (event) => {
    onOptionsChange({
      ...options,
      jsonData: { ...jsonData, timeout: event.target.value }
    });
  };
  return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(import_ui.InlineFieldRow, { children: [
    /* @__PURE__ */ (0, import_jsx_runtime.jsx)(
      import_ui.InlineField,
      {
        label: "Default interval",
        labelWidth: 16,
        tooltip: "Default ping interval when a query does not specify one, e.g. 10s",
        children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(
          import_ui.Input,
          {
            value: jsonData.defaultInterval || "",
            onChange: onIntervalChange,
            placeholder: "10s",
            width: 12
          }
        )
      }
    ),
    /* @__PURE__ */ (0, import_jsx_runtime.jsx)(
      import_ui.InlineField,
      {
        label: "Ping timeout",
        labelWidth: 16,
        tooltip: "Per-ping timeout, e.g. 3s",
        children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(import_ui.Input, { value: jsonData.timeout || "", onChange: onTimeoutChange, placeholder: "3s", width: 12 })
      }
    )
  ] });
}

// src/components/QueryEditor.tsx
var import_ui2 = require("@grafana/ui");
var import_jsx_runtime2 = require("react/jsx-runtime");
function QueryEditor({ query, onChange }) {
  const onHostChange = (event) => {
    onChange({ ...query, host: event.target.value });
  };
  const onIntervalChange = (event) => {
    onChange({ ...query, interval: event.target.value });
  };
  return /* @__PURE__ */ (0, import_jsx_runtime2.jsxs)(import_ui2.InlineFieldRow, { children: [
    /* @__PURE__ */ (0, import_jsx_runtime2.jsx)(import_ui2.InlineField, { label: "Host", labelWidth: 10, tooltip: "Hostname or IP address to ping, e.g. 8.8.8.8", children: /* @__PURE__ */ (0, import_jsx_runtime2.jsx)(import_ui2.Input, { value: query.host || "", onChange: onHostChange, placeholder: "8.8.8.8", width: 24 }) }),
    /* @__PURE__ */ (0, import_jsx_runtime2.jsx)(
      import_ui2.InlineField,
      {
        label: "Interval",
        labelWidth: 10,
        tooltip: "How often to ping this host, e.g. 10s, 30s, 1m",
        children: /* @__PURE__ */ (0, import_jsx_runtime2.jsx)(import_ui2.Input, { value: query.interval || "", onChange: onIntervalChange, placeholder: "10s", width: 12 })
      }
    )
  ] });
}

// src/module.ts
var plugin = new import_data2.DataSourcePlugin(DataSource).setConfigEditor(ConfigEditor).setQueryEditor(QueryEditor);
