truncate table datasets;

truncate table images;

truncate table labels;

truncate table project_permissions;

truncate table projects;

truncate table task_detail_0;

truncate table task_detail_1;

truncate table task_detail_2;

truncate table task_detail_3;

truncate table task_detail_4;

truncate table task_detail_5;

truncate table task_detail_6;

truncate table task_detail_7;

truncate table task_detail_8;

truncate table task_detail_9;

truncate table tasks;

truncate table workspace_permissions;

truncate table workspaces;

ALTER TABLE workspaces MODIFY description TEXT;
ALTER TABLE projects MODIFY description TEXT;
ALTER TABLE datasets MODIFY description TEXT;
ALTER TABLE tasks MODIFY description TEXT;

