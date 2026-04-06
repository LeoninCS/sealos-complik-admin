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
  const { configRecords, deleteConfigRecord } = useAppData();
  const [selected, setSelected] = useState<ConfigRecord | null>(null);
  const [open, setOpen] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [pendingDelete, setPendingDelete] = useState<ConfigRecord | null>(null);

  const rows = useMemo(() => {
    return configRecords.filter((item) => item.configName.toLowerCase().includes(keyword.toLowerCase()));
  }, [configRecords, keyword]);

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
        onClose={() => setOpen(false)}
        open={open}
        title="新增配置"
      >
        <div className="panel-stack">
          <Field label="配置名">
            <Input placeholder="例如：project-config-demo" />
          </Field>
          <Field label="配置类型">
            <Select defaultValue="json">
              <option value="json">json</option>
            </Select>
          </Field>
          <Field label="描述">
            <Input placeholder="简短说明用途" />
          </Field>
          <Field label="JSON 内容">
            <TextArea defaultValue={`{\n  "enabled": true,\n  "threshold": 3\n}`} />
          </Field>
          <div className="button-row">
            <Button variant="primary" onClick={() => setOpen(false)}>保存配置</Button>
            <Button variant="secondary" onClick={() => setOpen(false)}>取消</Button>
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
