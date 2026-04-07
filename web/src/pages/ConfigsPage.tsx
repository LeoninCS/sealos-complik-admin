import { useMemo, useState } from "react";
import {
  Button,
  ConfirmModal,
  DetailList,
  Drawer,
  EmptyState,
  Field,
  Input,
  Modal,
  PageHeader,
  Select,
  SurfaceCard,
  TextArea,
} from "../components/ui";
import { useAppData } from "../contexts/AppDataContext";
import type { ConfigRecord } from "../types";

export function ConfigsPage() {
  const { configRecords, createConfigRecord, deleteConfigRecord } = useAppData();
  const [selected, setSelected] = useState<ConfigRecord | null>(null);
  const [open, setOpen] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [pendingDelete, setPendingDelete] = useState<ConfigRecord | null>(null);
  const [configName, setConfigName] = useState("");
  const [configType, setConfigType] = useState("json");
  const [description, setDescription] = useState("");
  const [value, setValue] = useState('{\n  "enabled": true,\n  "threshold": 3\n}');
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const rows = useMemo(() => {
    return configRecords.filter((item) => item.configName.toLowerCase().includes(keyword.toLowerCase()));
  }, [configRecords, keyword]);

  const handleCreateConfig = async () => {
    if (!configName.trim() || !configType.trim() || !value.trim()) {
      setFormError("配置名、配置类型和 JSON 内容均为必填。");
      return;
    }

    try {
      JSON.parse(value);
    } catch {
      setFormError("JSON 内容格式不正确，请检查后再提交。");
      return;
    }

    setSubmitting(true);
    setFormError(null);
    try {
      await createConfigRecord({
        configName: configName.trim(),
        configType: configType.trim(),
        description: description.trim(),
        value: value.trim(),
      });
      setOpen(false);
      setConfigName("");
      setConfigType("json");
      setDescription("");
      setValue('{\n  "enabled": true,\n  "threshold": 3\n}');
    } catch (err) {
      setFormError(err instanceof Error ? err.message : "新增配置失败");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="page-container">
      <PageHeader
        kicker="Configs"
        title="项目配置"
        description="统一查看配置名、类型、描述和 JSON 内容，新增和编辑保持同一套表单结构。"
        actions={<Button variant="primary" onClick={() => setOpen(true)}>新增配置</Button>}
      />

      <SurfaceCard>
        <div className="toolbar">
          <Field label="配置名搜索">
            <Input placeholder="按 config_name 搜索" value={keyword} onChange={(event) => setKeyword(event.target.value)} />
          </Field>
          <Field label="配置类型">
            <Select defaultValue="all">
              <option value="all">全部类型</option>
              <option value="json">JSON</option>
            </Select>
          </Field>
        </div>
      </SurfaceCard>

      <SurfaceCard className="data-table-wrap" padded={false}>
        {rows.length > 0 ? (
          <table className="data-table">
            <thead>
              <tr>
                <th>配置名</th>
                <th>类型</th>
                <th>描述</th>
                <th>更新时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((item) => (
                <tr key={item.id}>
                  <td>
                    <button className="namespace-link table-row-button" onClick={() => setSelected(item)} type="button">
                      {item.configName}
                    </button>
                  </td>
                  <td>{item.configType}</td>
                  <td>{item.description}</td>
                  <td>{item.updatedAt}</td>
                  <td>
                    <Button variant="ghost" onClick={() => setSelected(item)}>
                      查看
                    </Button>
                    <Button variant="danger" onClick={() => setPendingDelete(item)}>
                      删除
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <div style={{ padding: 20 }}>
            <EmptyState
              title="还没有项目配置"
              description="当前筛选条件下没有配置记录，可以新建一条配置。"
              action={<Button variant="primary" onClick={() => setOpen(true)}>新增配置</Button>}
            />
          </div>
        )}
      </SurfaceCard>

      <Drawer
        description="右侧抽屉展示配置详情，并保留足够宽度展示 JSON 内容。"
        onClose={() => setSelected(null)}
        open={Boolean(selected)}
        title={selected ? selected.configName : ""}
      >
        {selected ? (
          <>
            <DetailList
              items={[
                { label: "配置名", value: selected.configName },
                { label: "配置类型", value: selected.configType },
                { label: "描述", value: selected.description },
                { label: "更新时间", value: selected.updatedAt },
              ]}
            />
            <div style={{ marginTop: 20 }}>
              <div className="detail-label" style={{ marginBottom: 8 }}>JSON 内容</div>
              <pre className="code-block">{selected.value}</pre>
            </div>
            <div className="button-row" style={{ marginTop: 20 }}>
              <Button variant="danger" onClick={() => setPendingDelete(selected)}>
                删除配置
              </Button>
            </div>
          </>
        ) : null}
      </Drawer>

      <Modal
        description="演示用表单，结构与后续真实接口表单保持一致。"
        onClose={() => {
          setOpen(false);
          setFormError(null);
        }}
        open={open}
        title="新增配置"
      >
        <div className="panel-stack">
          <Field label="配置名">
            <Input placeholder="例如：project-config-demo" value={configName} onChange={(event) => setConfigName(event.target.value)} />
          </Field>
          <Field label="配置类型">
            <Select value={configType} onChange={(event) => setConfigType(event.target.value)}>
              <option value="json">json</option>
            </Select>
          </Field>
          <Field label="描述">
            <Input placeholder="简短说明用途" value={description} onChange={(event) => setDescription(event.target.value)} />
          </Field>
          <Field label="JSON 内容">
            <TextArea value={value} onChange={(event) => setValue(event.target.value)} />
          </Field>
          {formError ? <div className="muted-text" style={{ color: "#b42318" }}>{formError}</div> : null}
          <div className="button-row">
            <Button variant="primary" onClick={() => void handleCreateConfig()}>
              {submitting ? "保存中..." : "保存配置"}
            </Button>
            <Button
              variant="secondary"
              onClick={() => {
                setOpen(false);
                setFormError(null);
              }}
            >
              取消
            </Button>
          </div>
        </div>
      </Modal>

      <ConfirmModal
        description={pendingDelete ? `删除后将从当前前端列表中移除配置 ${pendingDelete.configName}。` : ""}
        onClose={() => setPendingDelete(null)}
        onConfirm={() => {
          if (!pendingDelete) return;
          void deleteConfigRecord(pendingDelete.configName).then(() => {
            if (selected?.id === pendingDelete.id) {
              setSelected(null);
            }
            setPendingDelete(null);
          });
        }}
        open={Boolean(pendingDelete)}
        title="删除配置"
      />
    </div>
  );
}
