const BASE_URL = "http://localhost:8080/tasks/"

document.addEventListener('DOMContentLoaded', () => {

    loadTasksAndCategories();

    document.getElementById('category-form').addEventListener('submit', (e) => {
        e.preventDefault();
        const categoryTitle = document.getElementById('category-title').value;
        if (categoryTitle.trim() === '') return;

        let rqbody = {
            order: 1,
            name: categoryTitle
        }

        let request = new Request(BASE_URL + "addCategory", {
            body: JSON.stringify(rqbody),
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });
        fetch(request).then(response => {
            if (!response.ok) {
                throw new Error("Network resposne was not ok");
            }
            return response.json();
        }).then(category => {
            addCategory(category);
        }).catch(() => {
            alert("Fehler beim hinzufügen!");
        })
    });

    document.addEventListener('change', (e) => {
        if (e.target.classList.contains('editable')) {
            let state = 0;
            if (e.target.parentElement.style.textDecoration == "line-through") {
                state = 1;
            }
            let belongs_to = parseInt(e.target.parentElement.parentElement.id.replace("c", ""));
            let rqbody = {
                id: parseInt(e.target.parentElement.id.replace("t", "")),
                title: e.target.parentElement.querySelector(".task-title-input").value,
                details: e.target.parentElement.querySelector(".task-details-input").value,     
                due: e.target.parentElement.querySelector(".task-due-input").value,
                state: state,
                belongs_to: belongs_to
            }
    
            let request = new Request(BASE_URL + "updateTask", {
                body: JSON.stringify(rqbody),
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                }
            });
            fetch(request).then(response => {
                if (!response.ok) {
                    throw new Error("Network resposne was not ok");
                }
                return response.json();
            }).then(() => {
                return;
            }).catch(() => {
                alert("Fehler beim ändern!");
            })
        }
    });

    document.addEventListener('click', (e) => {
        if (e.target.classList.contains('editable')) {
            let state = 0;
            if (e.target.parentElement.style.textDecoration == "line-through") {
                state = 1;
            }
            let belongs_to = parseInt(e.target.parentElement.parentElement.id.replace("c", ""));
            let rqbody = {
                id: parseInt(e.target.parentElement.id.replace("t", "")),
                title: e.target.parentElement.querySelector(".task-title-input").value,
                details: e.target.parentElement.querySelector(".task-details-input").value,     
                due: e.target.parentElement.querySelector(".task-due-input").value,
                state: state,
                belongs_to: belongs_to
            }
    
            let request = new Request(BASE_URL + "updateTask", {
                body: JSON.stringify(rqbody),
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                }
            });
            fetch(request).then(response => {
                if (!response.ok) {
                    throw new Error("Network resposne was not ok");
                }
                return response.json();
            }).then(() => {
                return;
            }).catch(() => {
                alert("Fehler beim ändern!");
            })
        }
    });
    
    let draggedItem = null;

    document.addEventListener('dragstart', (e) => {
        if (e.target.classList.contains('draggable')) {
            draggedItem = e.target;
            setTimeout(() => e.target.classList.add('hidden'), 0);
        }
    });

    document.addEventListener('dragend', (e) => {
        if (draggedItem) {
            draggedItem.classList.remove('hidden');
            draggedItem = null;
        }
    });

    document.addEventListener('dragover', (e) => {
        e.preventDefault();
        let classString = null;
        if (draggedItem.classList.contains('todo-item')) {
            classString = '.todo-item:not(.hidden), .category-title';
        }
        if (draggedItem.classList.contains('category-container')) {
            classString = '.category-container:not(.hidden)';
        }
        const closestTaskBottom = getClosestBottom(e.clientY, classString);
        const closestTaskTop = getClosestTop(e.clientY, classString);
        if (closestTaskBottom && closestTaskTop) {
            let bottombox = closestTaskBottom.getBoundingClientRect();
            let closestTaskBottomOffset = Math.abs(e.clientY - bottombox.top - bottombox.height / 2);
            let topbox = closestTaskTop.getBoundingClientRect();
            let closestTaskTopOffset = e.clientY - topbox.top - topbox.height / 2;
            if (closestTaskBottomOffset > closestTaskTopOffset) {
                if (closestTaskTopOffset < 50) {
                    e.preventDefault();
                    if ((closestTaskTop.classList.contains('category-title')) && (closestTaskTop.parentElement.parentElement.querySelectorAll('.todo-item').length == 0)) {
                        closestTaskTop.parentElement.parentElement.appendChild(draggedItem);
                    } else if (!(closestTaskTop.classList.contains('category-title'))) {
                    closestTaskTop.after(draggedItem);
                    }
                }
            } else {
                if (closestTaskBottomOffset < 50) {
                    e.preventDefault();
                    if ((closestTaskBottom.classList.contains('category-title')) && (closestTaskBottom.parentElement.parentElement.querySelectorAll('.todo-item').length == 0)) {
                        closestTaskBottom.parentElement.parentElement.appendChild(draggedItem);
                    } else if (!(closestTaskBottom.classList.contains('category-title'))) {
                    closestTaskBottom.before(draggedItem);
                    }
                }
            }
        } else if (closestTaskTop) {
            e.preventDefault();
            if ((closestTaskTop.classList.contains('category-title')) && (closestTaskTop.parentElement.parentElement.querySelectorAll('.todo-item').length == 0)) {
                closestTaskTop.parentElement.parentElement.appendChild(draggedItem);
            } else if (!(closestTaskTop.classList.contains('category-title'))) {
            closestTaskTop.after(draggedItem);
            }
        } else if (closestTaskBottom) {
            e.preventDefault();
            if ((closestTaskTop.classList.contains('category-title')) && (closestTaskTop.parentElement.parentElement.querySelectorAll('.todo-item').length == 0)) {
                closestTaskTop.parentElement.parentElement.appendChild(draggedItem);
            } else if (!(closestTaskBottom.classList.contains('category-title'))) {
            closestTaskTop.before(draggedItem);
            }
        }
        
    });

    document.addEventListener('drop', (e) => {
        if (draggedItem) {
            e.preventDefault();
            switch (draggedItem.className) {
                case 'category-container draggable hidden':
                    let categoryArray = [...document.querySelectorAll('.category-container')];
                    let rqbody = {
                        id: parseInt(draggedItem.id.replace("c", "")),
                        order: categoryArray.indexOf(draggedItem) + 1
                    }
                    let request = new Request(BASE_URL + "relocateCategory", {
                        body: JSON.stringify(rqbody),
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json"
                        }
                    });
                    fetch(request).then(response => {
                        if (!response.ok) {
                            throw new Error("Network resposne was not ok");
                        }
                        return response.json();
                    }).then(() => {
                        location.reload();
                    }).catch(() => {
                        alert("Fehler beim ändern!");
                    })
                break;
                case 'todo-item draggable hidden':
                    let state = 0;
                    if (draggedItem.style.textDecoration == "line-through") {
                        state = 1;
                    }
                    let taskArray = Array.from(draggedItem.parentElement.querySelectorAll('.todo-item'));
                    let rqbodytask = {
                        id: parseInt(draggedItem.id.replace("t", "")),
                        title: draggedItem.querySelector(".task-title-input").value,
                        details: draggedItem.querySelector(".task-details-input").value,     
                        due: draggedItem.querySelector(".task-due-input").value,
                        state: state,
                        belongs_to: parseInt(draggedItem.parentElement.id.replace("c", "")),
                        order: taskArray.indexOf(draggedItem) + 1
                    }
                    let requesttask = new Request(BASE_URL + "relocateTask", {
                        body: JSON.stringify(rqbodytask),
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json"
                        }
                    });
                    fetch(requesttask).then(response => {
                        if (!response.ok) {
                            throw new Error("Network resposne was not ok");
                        }
                        return response.json();
                    }).then((taskdata) => {
                        location.reload();
                    }).catch(() => {
                        alert("Fehler beim ändern!");
                    })
                break;
            }
            
        }
    })

    function getClosestBottom(y, classString) {
        const tasks = [...document.querySelectorAll(classString)];
        return tasks.reduce((closest, task) => {
            const box = task.getBoundingClientRect();
            const offset = y - box.top - box.height / 2;
            const condition = offset < 0 && offset > closest.offset;
            if (condition) {
                return { offset: offset, element: task };
            } else {
                return closest;
            }
        }, { offset: Number.NEGATIVE_INFINITY }).element;
    }

    function getClosestTop(y, classString) {
        const tasks = [...document.querySelectorAll(classString)];
        return tasks.reduce((closest, task) => {
            const box = task.getBoundingClientRect();
            const offset = y - box.top - box.height / 2;
            const condition = offset > 0 && offset < closest.offset;
            if (condition) {
                return { offset: offset, element: task };
            } else {
                return closest;
            }
        }, { offset: Number.POSITIVE_INFINITY }).element;
    }
});

function addCategory(categorydata) {
    

    const categoryContainer = document.createElement('div');
    categoryContainer.className = 'category-container draggable';
    categoryContainer.id = "c" + categorydata.id;
    categoryContainer.draggable = true;

    const categoryHeader = document.createElement('div');
    categoryHeader.className = 'category-header';

    const categoryTitleField = document.createElement('input');
    categoryTitleField.className = 'category-title';
    categoryTitleField.value = categorydata.name;
    categoryTitleField.onchange = () => {
        
        let rqbody = {
            id: categorydata.id,
            name: categoryTitleField.value
        }

        let request = new Request(BASE_URL + "updateCategory", {
            body: JSON.stringify(rqbody),
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });
        fetch(request).then(response => {
            if (!response.ok) {
                throw new Error("Network resposne was not ok");
            }
            return response.json();
        }).then(category => {
            categoryTitleField.value = category.name;
        }).catch(() => {
            alert("Fehler beim ändern!");
        })
    }

    const deleteCategoryButton = document.createElement('button');
    deleteCategoryButton.className = 'category-delete';
    deleteCategoryButton.textContent = 'Kategorie löschen';
    deleteCategoryButton.onclick = () => {
        let request = new Request(BASE_URL + "deleteCategory", {
            body: JSON.stringify({id: categorydata.id}),
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });
        fetch(request)
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
        })
        .then(() => {
            categoryContainer.remove();
        })
        .catch(error => {
            console.error('Error:', error);
            alert("Fehler beim Löschen der Kategorie!");
        });
        }

    categoryHeader.appendChild(categoryTitleField);
    categoryHeader.appendChild(deleteCategoryButton);
    categoryContainer.appendChild(categoryHeader);

    const addButton = document.createElement('button');
    addButton.textContent = 'Aufgabe hinzufügen';
    addButton.className = 'add-button';
    addButton.onclick = () => {
        let request = new Request(BASE_URL + "addTask", {
            body: JSON.stringify({
                belongs_to: categorydata.id,
                title: "",
                order: 1,
                state: 0,
                due: "",
                details: ""
            }),
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });
        fetch(request)
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            } else {
                return response.json();
            }
        })
        .then((taskdata) => {
            addTodoItem(categoryContainer, taskdata);
        })
        .catch(error => {
            console.error('Error:', error);
            alert("Fehler beim Hinzufügen der Aufgabe!");
        });
    }
    categoryContainer.appendChild(addButton);

    document.getElementById('categories').append(categoryContainer);

    document.getElementById('category-title').value = '';

    return categoryContainer;
}

function addTodoItem(categoryContainer, taskdata) {
    const todoItem = document.createElement('div');
    todoItem.className = 'todo-item draggable';
    todoItem.draggable = true;
    todoItem.id = "t" + taskdata.id;

    const dragHandle = document.createElement('span');
    dragHandle.className = 'drag-handle';
    dragHandle.innerHTML = '<i class="fas fa-grip-lines"></i>';

    const titleInput = document.createElement('input');
    titleInput.type = 'text';
    titleInput.placeholder = 'Titel';
    titleInput.className = 'editable task-title-input';
    titleInput.value = taskdata.title;

    const detailsInput = document.createElement('input');
    detailsInput.type = 'text';
    detailsInput.placeholder = 'Details';
    detailsInput.className = 'editable task-details-input';
    detailsInput.value = taskdata.details;
    
    const dueDateInput = document.createElement('input');
    dueDateInput.type = 'datetime-local';
    dueDateInput.className = 'editable task-due-input';
    dueDateInput.value = taskdata.due;

    const completeButton = document.createElement('button');
    completeButton.textContent = 'Erledigt';
    completeButton.className = 'complete editable';

    if (taskdata.state != 0) {
        todoItem.style.textDecoration = 'line-through';
    }
    completeButton.onclick = () => {
        if (todoItem.style.textDecoration === 'line-through') {
            todoItem.style.textDecoration = 'none';
        } else {
            todoItem.style.textDecoration = 'line-through';
        }
    };

    const deleteButton = document.createElement('button');
    deleteButton.textContent = 'Löschen';
    deleteButton.className = 'delete';
    deleteButton.onclick = () => {
        let taskArray = Array.from(deleteButton.parentElement.parentElement.querySelectorAll('.todo-item'));
        let rqbody = {
            belongs_to: taskdata.belongs_to,
            id: taskdata.id,
            order: taskArray.indexOf(deleteButton.parentElement) + 1
        };
        let request = new Request(BASE_URL + "deleteTask", {
            body: JSON.stringify(rqbody),
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });
        fetch(request)
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
        })
        .then(() => {
            categoryContainer.removeChild(todoItem);
        })
        .catch(error => {
            console.error('Error:', error);
            alert("Fehler beim Hinzufügen der Aufgabe!");
        });
        
    }

    todoItem.appendChild(dragHandle);
    todoItem.appendChild(titleInput);
    todoItem.appendChild(detailsInput);
    todoItem.appendChild(dueDateInput);
    todoItem.appendChild(completeButton);
    todoItem.appendChild(deleteButton);

    categoryContainer.appendChild(todoItem);
}

function renderTodoList(data) {
    const container = document.getElementById('categories');
    container.innerHTML = '';

    if (data.categories == null) {
        return;
    }

    // Sort categories by their order
    const sortedCategories = data.categories.sort((a, b) => a.order - b.order);

    sortedCategories.forEach(category => {
        // Create category container
        let container = addCategory(category);

        // Filter tasks by category id and sort by order
        if (data.tasks != null) {
        const categoryTasks = data.tasks
            .filter(task => task.belongs_to === category.id)
            .sort((a, b) => a.order - b.order);

        categoryTasks.forEach(task => {
            addTodoItem(container, task);
        });}
    });
}

function loadTasksAndCategories() {
    const URL = "http://localhost:8080/tasks/get";
    
    fetch(URL)
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then(data => {
            renderTodoList(data);
        })
        .catch(error => {
            console.error('Error:', error);
            alert("Fehler beim Laden der Aufgaben!");
        });
}
