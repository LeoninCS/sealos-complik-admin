import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
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
} from "../components/ui";
import { useAppData } from "../contexts/AppDataContext";
import type { UnbanRecord } from "../types";

export function UnbansPage() {
  const navigate = useNavigate();
  const { unbanRecords, deleteUnbanRecord } = useAppData();
  const [selected, setSelected] = useState<UnbanRecord | null>(null);
  const [open, setOpen] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [pendingDelete, setPendingDelete] = useState<UnbanRecord | null>(null);

  const rows = useMemo(() => {
    return unbanRecords.filter((item) => item.namespace.toLowerCase().includes(keyword.toLowerCase()));
  }, [keyword, unbanRecords]);

  return (
    <div className="page-container">
      <PageHeader
        kicker="Unbans"
        title="解封记录"
        description="页面复杂度低于封禁记录，只保留必要筛选和最小录入字段。"
        actions={<Button variant="primary" onClick={() => setOpen(true)}>新增解封</Button>}
      />

      <SurfaceCard>
        <div className="toolbar">
          <Field label="namespace">
            <Input placeholder="按 namespace 搜索" value={keyword} onChange={(event) => setKeyword(event.target.value)} />
          </Field>
          <Field label="操作人">
            <Input placeholder="例如：Bob" />
          </Field>
          <Field label="时间范围">
            <Select defaultValue="7d">
              <option value="24h">最近 24 小时</option>
              <option value="7d">最近 7 天</option>
              <option value="30d">最近 30 天</option>
            </Select>
          </Field>
        </div>
      </SurfaceCard>

      <SurfaceCard className="data-table-wrap" padded={false}>
        {rows.length > 0 ? (
          <table className="data-table">
            <thead>
              <tr>
                <th>namespace</th>
                <th>操作人</th>
                <th>时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((item) => (
                <tr key={item.id}>
                  <td>
                    <button className="namespace-link table-row-button" onClick={() => navigate(`/namespaces/${item.namespace}`)} type="button">
                      {item.namespace}
                    </button>
                  </td>
                  <td>
                    <button className="table-row-button" onClick={() => setSelected(item)} type="button">
                      {item.operatorName}
                    </button>
                  </td>
                  <td>{item.createdAt}</td>
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
              title="当前没有解封记录"
              description="可以直接新增一条解封记录。"
              action={<Button variant="primary" onClick={() => setOpen(true)}>新增解封</Button>}
            />
          </div>
        )}
      </SurfaceCard>

      <Drawer
        description="解封记录详情保持轻量，只展示当前接口已有字段。"
        onClose={() => setSelected(null)}
        open={Boolean(selected)}
        title={selected ? selected.namespace : ""}
      >
        {selected ? (
          <>
            <DetailList
              items={[
                { label: "namespace", value: selected.namespace },
                { label: "操作人", value: selected.operatorName },
                { label: "创建时间", value: selected.createdAt },
              ]}
            />
            <div className="button-row" style={{ marginTop: 20 }}>
              <Button variant="secondary" onClick={() => navigate(`/namespaces/${selected.namespace}`)}>
                查看 namespace 详情
              </Button>
              <Button variant="danger" onClick={() => setPendingDelete(selected)}>
                删除记录
              </Button>
            </div>
          </>
        ) : null}
      </Drawer>

      <Modal
        description="只保留 namespace 和操作人两个必填项。"
        onClose={() => setOpen(false)}
        open={open}
        title="新增解封"
      >
        <div className="panel-stack">
          <Field label="namespace">
            <Input placeholder="例如：growth-ops" />
          </Field>
          <Field label="操作人">
            <Input placeholder="例如：Bob" />
          </Field>
          <div className="button-row">
            <Button variant="primary" onClick={() => setOpen(false)}>保存解封记录</Button>
            <Button variant="secondary" onClick={() => setOpen(false)}>取消</Button>
          </div>
        </div>
      </Modal>

      <ConfirmModal
        description={pendingDelete ? `删除后将从当前前端列表中移除 namespace ${pendingDelete.namespace} 的解封记录。` : ""}
        onClose={() => setPendingDelete(null)}
        onConfirm={() => {
          if (!pendingDelete) return;
          void deleteUnbanRecord(pendingDelete.namespace).then(() => {
            if (selected?.id === pendingDelete.id) {
              setSelected(null);
            }
            setPendingDelete(null);
          });
        }}
        open={Boolean(pendingDelete)}
        title="删除解封记录"
      />
    </div>
  );
}
