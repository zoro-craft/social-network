# Groups

## Overview

The **Groups** feature allows users to create and participate in private communities inside the social network.

A group has a **title** and a **description**. Users can join a group either by being invited by an existing member or by requesting to join a group they found through the group browser.

Only members of a group can see and interact with its content.

---

## 1. Creating a Group

A user can create a new group.

When creating a group, the user must provide:

* **Title**
* **Description**

The user who creates the group becomes the **creator** of the group.

The creator has special permissions, especially for managing join requests.

---

## 2. Inviting Users

The creator can invite other users to join the group.

Users who are already members of the group can also invite other users.

An invitation does **not** automatically make the user a member.

The invited user must:

* **Accept** the invitation → becomes a group member.
* **Reject** the invitation → does not join the group.

### Invitation flow

```text
Member
   │
   ├── Invite User
   │
   ▼
Invitation
   │
   ├── Accept ──► User becomes a Member
   │
   └── Reject ──► Invitation refused
```

---

## 3. Requesting to Join a Group

Users can also join a group without receiving an invitation.

To do this, the user must first find the group through the **Groups** section.

The user can then send a **join request**.

Only the **creator of the group** can manage these requests.

The creator can:

* Accept the request → user becomes a member.
* Refuse the request → user does not join the group.

### Join request flow

```text
User
  │
  ▼
Browse Groups
  │
  ▼
Select Group
  │
  ▼
Request to Join
  │
  ▼
Group Creator
  │
  ├── Accept ──► User becomes a Member
  │
  └── Refuse ──► Request rejected
```

---

## 4. Browsing Groups

There must be a section where users can browse through the available groups.

This allows a user to:

1. View available groups.
2. Find a group they are interested in.
3. Open the group's information.
4. Request to join the group.

A group should display at least its:

* Title
* Description

---

## 5. Group Posts

Once a user becomes a member of a group, they can create posts inside that group.

Group members can also comment on posts.

### Important Rule

Group posts and comments are **private to the group**.

Only members of that group can:

* See group posts.
* Create posts.
* See comments.
* Create comments.

Users who are not members must not be able to access the group's posts or comments.

```text
Group
 │
 ├── Post 1
 │    ├── Comment
 │    └── Comment
 │
 ├── Post 2
 │    └── Comment
 │
 └── Post 3

      ↓

Only Group Members can see them
```

---

## 6. Group Events

A member of a group can create an **event** for the other members.

An event must contain:

* **Title**
* **Description**
* **Day / Time**
* At least two participation options:

  * **Going**
  * **Not going**

After the event is created, every group member can choose one of the available options.

### Example

```text
Event: Football Match

Description:
Football match between Zone01 students.

Date:
Saturday - 18:00

Options:
[ Going ]
[ Not going ]
```

Each member can select their participation status.

---

## 7. Group Permissions

The main permissions can be summarized as follows:

| Action                 | Creator | Member | Non-member |
| ---------------------- | ------: | -----: | ---------: |
| Create group           |       ✅ |      ✅ |          ✅ |
| Invite users           |       ✅ |      ✅ |          ❌ |
| Accept invitations     |       ✅ |      ✅ |          ❌ |
| Request to join        |       ❌ |      ❌ |          ✅ |
| Accept join requests   |       ✅ |      ❌ |          ❌ |
| Refuse join requests   |       ✅ |      ❌ |          ❌ |
| Create group posts     |       ✅ |      ✅ |          ❌ |
| Comment on group posts |       ✅ |      ✅ |          ❌ |
| See group posts        |       ✅ |      ✅ |          ❌ |
| Create events          |       ✅ |      ✅ |          ❌ |
| Participate in events  |       ✅ |      ✅ |          ❌ |

> The creator is also a member of the group.

---

## 8. Complete Group Workflow

The complete feature can be represented as:

```text
                    ┌─────────────────┐
                    │  Create Group   │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Group Creator   │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
              ▼                             ▼
      Invite Users                  Browse Groups
              │                             │
              ▼                             ▼
        Invitation                  Join Request
              │                             │
        ┌─────┴─────┐                       │
        │           │                       ▼
      Accept      Reject              Group Creator
        │                                   │
        ▼                              ┌────┴────┐
     Member                            │         │
                                      Accept    Refuse
                                        │
                                        ▼
                                      Member
```

After becoming a member:

```text
                 Group Member
                      │
          ┌───────────┼───────────┐
          │           │           │
          ▼           ▼           ▼
       Create       Comment      Create
        Post        Posts        Event
          │           │           │
          └───────────┴───────────┘
                      │
                      ▼
              Other Group Members
                      │
                      ▼
              Participate in Event
              ┌───────────────┐
              │ Going         │
              │ Not going     │
              └───────────────┘
```

---

## 9. Main Requirements

The Groups feature must provide:

* [ ] Create a group with a title and description.
* [ ] Identify the creator of each group.
* [ ] Browse available groups.
* [ ] Invite users to a group.
* [ ] Accept or reject invitations.
* [ ] Allow group members to invite other users.
* [ ] Request to join a group.
* [ ] Allow only the creator to accept or refuse join requests.
* [ ] Allow members to create group posts.
* [ ] Allow members to comment on group posts.
* [ ] Make group posts and comments visible only to members.
* [ ] Allow members to create events.
* [ ] Events must contain title, description, and date/time.
* [ ] Events must provide at least `Going` and `Not going` options.
* [ ] Allow group members to select their event participation status.

---

## 10. Important Concept

The most important rule of this feature is **group membership**.

A user should only have access to the group's private content and actions if they are a member.

```text
                    Is User a Member?
                           │
                    ┌──────┴──────┐
                    │             │
                   YES            NO
                    │             │
                    ▼             ▼
              Group Content    No Access
              & Actions
```

Therefore, membership should be checked both in the **frontend** for the user interface and, more importantly, in the **backend** to prevent unauthorized access.
